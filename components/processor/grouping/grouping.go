package grouping

import (
	"context"
	"fmt"

	//"slices"
	"time"

	redisv9 "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/lang"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/mq"
	"github.com/japannext/snooze/internal/opensearch/format"
	"github.com/japannext/snooze/internal/pipeline"
	"github.com/japannext/snooze/internal/redis"
	"github.com/japannext/snooze/internal/utils"
)

type Processor struct {
	groupings []Grouping
	storeQ    *mq.Pub
}

func New(cfgs []pipeline.Grouping) (*Processor, error) {
	p := &Processor{}

	p.storeQ = mq.StorePub()

	for _, cfg := range cfgs {
		gr, err := NewGrouping(cfg)
		if err != nil {
			return &Processor{}, fmt.Errorf("error loading group '%s': %w", cfg.Name, err)
		}

		p.groupings = append(p.groupings, gr)
	}

	return p, nil
}

// Group logs by fields. Groups can then be used
// for rate-limiting, snooze, and UI search.
type Grouping struct {
	pipeline.Grouping
	condition *lang.Condition
	fields    []*lang.Field
}

func NewGrouping(cfg pipeline.Grouping) (Grouping, error) {
	group := Grouping{Grouping: cfg}
	if group.If != "" {
		condition, err := lang.NewCondition(group.If)
		if err != nil {
			return Grouping{}, fmt.Errorf("error loading condition `%s`: %w", group.If, err)
		}

		group.condition = condition
	}

	if len(group.GroupBy) > 0 {
		fields, err := lang.NewFields(group.GroupBy)
		if err != nil {
			return Grouping{}, fmt.Errorf("error loading fields: %w", err)
		}

		group.fields = fields
	}

	if group.GroupByMap != "" && len(group.GroupBy) != 0 {
		return Grouping{}, fmt.Errorf("group_by and group_by_map are mutually exclusive")
	}

	return group, nil
}

const expiration = time.Duration(6) * time.Hour

func gkey(gr *models.Group) string {
	return fmt.Sprintf("grouping:%s:%s", gr.Name, gr.Hash)
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "grouping")
	defer span.End()

	for _, group := range p.groupings {
		if group.condition != nil {
			match, err := group.condition.MatchLog(ctx, item)
			if err != nil {
				log.Warnf("error while matching `%s`: %s", group.If, err)

				continue
			}

			if !match {
				continue
			}
		}

		gr := &models.Group{Name: group.Name, Labels: make(map[string]string), Score: group.Score}

		if len(group.GroupBy) > 0 {
			for _, field := range group.fields {
				value, err := lang.ExtractField(item, field)
				if err != nil {
					log.Warnf("Failed to match %s: %s", field, err)

					continue
				}

				gr.Labels[field.String()] = value
			}
		} else if group.GroupByMap != "" {
			switch group.GroupByMap {
			case "source":
				gr.Labels["source.kind"] = item.Source.Kind
				gr.Labels["source.name"] = item.Source.Name
			case "identity":
				for k, v := range item.Identity {
					gr.Labels["identity."+k] = v
				}
			case "labels":
				for k, v := range item.Labels {
					gr.Labels["labels."+k] = v
				}
			}
		}

		gr.Hash = utils.ComputeHash(gr.Labels)
		item.Groups = append(item.Groups, gr)
	}

	pipe := redis.Client.Pipeline()
	exists := make(map[string]*redisv9.IntCmd)

	for _, gr := range item.Groups {
		key := gkey(gr)
		exists[key] = pipe.Exists(ctx, key)
		// Set afterward anyway
		pipe.Set(ctx, key, "1", expiration)
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Errorf("failed to get/set groups from redis: %s", err)

		return decision.Retry(err)
	}

	// Insert the groups in Opensearch only if it doesn't exists
	for _, gr := range item.Groups {
		if cmd, ok := exists[gkey(gr)]; ok && cmd.Val() > 0 {
			// Skip insert due to LRU cache
			continue
		}

		err := p.storeQ.PublishData(ctx, &format.Update{
			Index:       models.GroupIndex,
			ID:          fmt.Sprintf("%s.%s", gr.Name, gr.Hash),
			Doc:         gr,
			DocAsUpsert: true,
		})
		if err != nil {
			log.Errorf("failed to publish group: %+v", gr)

			return decision.Retry(err)
		}
	}

	return decision.OK()
}

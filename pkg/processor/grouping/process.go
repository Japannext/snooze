package grouping

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	redisv9 "github.com/redis/go-redis/v9"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var expiration = time.Duration(6) * time.Hour

func gkey(gr *models.Group) string {
	return fmt.Sprintf("grouping:%s:%s", gr.Name, gr.Hash)
}

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "grouping")
	defer span.End()

	for _, group := range groupings {
		if group.internal.condition != nil {
			match, err := group.internal.condition.MatchLog(ctx, item)
			if err != nil {
				return err
			}
			if !match {
				continue
			}
		}

		var gr = &models.Group{Name: group.Name, Labels: make(map[string]string)}
		if len(group.GroupBy) > 0 {
			for _, field := range group.internal.fields {
				value, err := lang.ExtractField(item, field)
				if err != nil {
					logrus.Warnf("Failed to match %s: %s", field, err)
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
						gr.Labels[fmt.Sprintf("identity.%s", k)] = v
					}
				case "gr.Labels":
					for k, v := range item.Labels {
						gr.Labels[fmt.Sprintf("gr.Labels.%s", k)] = v
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
		log.Warnf("failed to get/set groups from redis: %s", err)
		tracing.Error(span, err)
		return err
	}

	// Insert the groups in Opensearch only if it doesn't exists
	for _, gr := range item.Groups {
		if cmd, ok := exists[gkey(gr)]; ok && cmd.Val() > 0 {
			// Skip insert due to LRU cache
			continue
		}

		err := storeQ.PublishData(ctx, &format.Index{
			Index: models.GROUP_INDEX,
			ID: gr.ID,
			Item: gr,
		})
		if err != nil {
			log.Warnf("failed to publish group: %+v", gr)
			continue
		}
	}

	return nil
}

package grouping

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/lang"
	"github.com/japannext/snooze/pkg/common/utils"
)

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

	return nil
}

package profile

import (
	"context"
	// "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/common/tracing"
)

var tracer = tracing.Tracer("snooze-process")

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracer.Start(ctx, "profile")
	defer span.End()
	for _, rule := range fastMapper.GetMatches(item) {
		reject := rule.Process(ctx, item)
		if reject {
			return nil
		}
	}

	return nil
}

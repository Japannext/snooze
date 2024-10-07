package profile

import (
	"context"
	// "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/tracing"
)

func Process(ctx context.Context, item *models.Log) error {
	ctx, span := tracing.TRACER.Start(ctx, "profile")
	defer span.End()
	for _, rule := range fastMapper.GetMatches(item) {
		reject := rule.Process(ctx, item)
		if reject {
			return nil
		}
	}

	return nil
}

package transform

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
	"github.com/japannext/snooze/pkg/processor/decision"
	"go.opentelemetry.io/otel"
	log "github.com/sirupsen/logrus"
)

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "transform")
	defer span.End()

	for _, tr := range p.transforms {
		if tr.condition != nil {
			match, err := tr.condition.MatchLog(ctx, item)
			if err != nil {
				log.Warnf("error while matching `%s`: %s", tr.cfg.If, err)

				continue
			}

			if !match {
				continue
			}
		}

		ctx = context.WithValue(ctx, "capture", map[string]string{})

		for index, action := range tr.actions {
			var err error

			ctx, err = action.Process(ctx, item)
			if err != nil {
				log.Warnf("error in transform %s#%d: %s", tr.cfg.Name, index+1, err)

				continue
			}
		}
	}

	return decision.OK()
}

package ratelimit

import (
	"github.com/japannext/snooze/pkg/common/mq"
	"github.com/japannext/snooze/pkg/common/utils"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var (
	log    *logrus.Entry
	tracer trace.Tracer
	rates  = []*RateLimit{}
	storeQ *mq.Pub
)

func Startup(ratelimits []*RateLimit) {
	initMetrics()
	log = logrus.WithFields(logrus.Fields{"module": "ratelimit"})
	tracer = otel.Tracer("snooze")
	storeQ = mq.StorePub()

	validator := utils.NewNameValidator(true)

	for _, rate := range ratelimits {
		if err := validator.Check(rate.Name); err != nil {
			log.Fatalf("error in ratelimits: %s", err)
		}
		rate.Load()
		rates = append(rates, rate)
	}
}

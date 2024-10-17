package ratelimit

import (
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/japannext/snooze/pkg/common/utils"
)

var log *logrus.Entry
var tracer trace.Tracer
var rates = []*RateLimit{}

func Startup(ratelimits []*RateLimit) {
	initMetrics()
	log = logrus.WithFields(logrus.Fields{"module": "ratelimit"})
	tracer = otel.Tracer("snooze")

	validator := utils.NewNameValidator(true)

	for _, rate := range ratelimits {
		if err := validator.Check(rate.Name); err != nil {
			log.Fatalf("error in ratelimits: %s", err)
		}
		rate.Load()
		rates = append(rates, rate)
	}
}

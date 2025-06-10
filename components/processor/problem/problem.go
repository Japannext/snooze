package problem

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pgsql"
	"github.com/japannext/snooze/internal/redis"
	"github.com/japannext/snooze/internal/tracing"
)

const ttl = 15 * time.Minute

type Processor struct {
}

func New() (*Processor, error) {
	return &Processor{}, nil
}

func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {
	ctx, span := otel.Tracer("snooze").Start(ctx, "problem")
	defer span.End()

	if item.Service != "" {
		if err := ensureProblem(ctx, item); err != nil {
			return decision.Retry(err)
		}
	}

	return decision.OK()
}

// Make sure a problem object is created in the database. If a problem
// object exists, continue.
func ensureProblem(ctx context.Context, item *models.Log) error {
	ctx, span := otel.Tracer("snooze").Start(ctx, "ensureProblem")
	defer span.End()

	// Ensure the problem ID is set/reserved for a given period of time.
	id, found, err := redis.ReserveID(ctx, redis.ServiceProblemKey(item.Service), ttl)
	if err != nil {
		return fmt.Errorf("failed to reserve ID: %w", err)
	}

	// Assign the found problemID to the log.
	item.ProblemID = id
	tracing.SetString(span, "problem.id", id)
	tracing.SetBool(span, "problem.found", found)

	// Create
	if !found {
		_, err := pgsql.Insert(ctx, &models.Problem{
			UUID:      id,
			Name:      "",
			Kind:      "problem",
			Status:    "active",
			StartTime: time.Now(),
		})
		if err != nil {
			return fmt.Errorf("failed to insert to pgsql: %w", err)
		}
	}

	return nil
}

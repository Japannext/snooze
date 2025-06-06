package problem

import (
	"context"
	"fmt"
	"time"

	"github.com/japannext/snooze/components/processor/decision"
	"github.com/japannext/snooze/internal/models"
	"github.com/japannext/snooze/internal/pgsql"
	"github.com/japannext/snooze/internal/redis"
)

const ttl = 15 * time.Minute

type Processor struct {
}

func New() (*Processor, error) {
	return nil, nil
}

// Create a problem if necessary
func (p *Processor) Process(ctx context.Context, item *models.Log) *decision.Decision {

	// Ensure the problem ID is set/reserved for a given period of time.
	id, found, err := redis.ReserveID(ctx, redis.ServiceProblemKey(item.Service), ttl)
	if err != nil {
		return decision.Retry(fmt.Errorf("failed to reserve ID: %w", err))
	}

	// Assign the found problemID to the log.
	item.ProblemID = id

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
			return decision.Retry(fmt.Errorf("failed to insert to pgsql: %w", err))
		}
	}

	return decision.OK()
}

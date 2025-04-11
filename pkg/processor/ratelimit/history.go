package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/redis"
	"github.com/japannext/snooze/pkg/models"
)

// Create a new ratelimit record in the database
func (p *Processor) newRecord(ctx context.Context, rule Ratelimit, gr *models.Group, throttle *redis.ThrottleState) (string, error) {
	item := &models.RatelimitRecord{
		StartsAt: uint64(time.Now().UnixMilli()),
		Rule:     rule.Name,
		Hash:     gr.Hash,
	}

	id, err := opensearch.Index(ctx, &opensearch.IndexReq{
		Index: models.RatelimitHistoryIndex,
		Item:  item,
	})
	if err != nil {
		return "", fmt.Errorf("failed to publish ratelimit record: %w", err)
	}

	return id, nil
}

// Close the ratelimit record in the database (set the end time, and set as inactive)
func (p *Processor) closeRecord(ctx context.Context, id string) error {
	err := opensearch.Update(ctx, opensearch.UpdateReq{
		Index: models.RatelimitHistoryIndex,
		ID:    id,
		Doc: map[string]interface{}{
			"active": false,
			"endsAt": int64(time.Now().UnixMilli()),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to close ratelimit record (id=%s): %w", id, err)
	}

	return nil
}

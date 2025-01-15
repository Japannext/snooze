package alertmanager

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/japannext/snooze/pkg/alertmanager/status"
	"github.com/japannext/snooze/pkg/common/opensearch"
	"github.com/japannext/snooze/pkg/common/opensearch/format"
	"github.com/japannext/snooze/pkg/models"
)

func closeExpiredAlerts() error {
	ctx := context.TODO()

	// Get list of expired alerts (alerts which score was not updated)
	expiredItems, err := status.GetExpired(ctx)
	if err != nil {
		return fmt.Errorf("failed to get list of expired alerts: %w", err)
	}

	ids := []string{}

	for _, expiredItem := range expiredItems {
		ids = append(ids, expiredItem.ID)
	}

	// Fetch the alaerts from opensearch
	req := &opensearch.SearchReq{Index: models.ActiveAlertIndex}
	req.Doc.WithTerms("_id", ids)

	list, err := opensearch.Search[*models.ActiveAlert](ctx, req)
	if err != nil {
		return fmt.Errorf("failed to search in opensearch: %w", err)
	}

	now := models.TimeNow()

	for _, item := range list.Items {
		// Add the item to history
		err := storeQ.PublishData(ctx, &format.Create{
			Index: models.AlertHistoryIndex,
			Item: models.AlertRecord{AlertBase: item.AlertBase, EndsAt: now},
		})
		if err != nil {
			log.Errorf("failed to publish alert to history: %s", err)

			continue
		}

		// Remove from list of active alerts
		err = storeQ.PublishData(ctx, &format.Delete{
			Index: models.ActiveAlertIndex,
			ID: item.ID,
		})
		if err != nil {
			log.Errorf("failed to remove alert from active alerts: %s", err)

			continue
		}

		// Remove from sorted set
		key := fmt.Sprintf("alertmanager/%s/%s/%s", item.AlertGroup, item.AlertName, item.Hash)
		if err := status.Delete(ctx, key); err != nil {
			log.Errorf("failed to dequeue the alert from sorted set: %s", err)

			continue
		}
	}

	return nil
}

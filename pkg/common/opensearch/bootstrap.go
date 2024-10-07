package opensearch

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func bootstrap(ctx context.Context) {
	for _, tpl := range models.INDEXES {
		ensureIndex(ctx, tpl.IndexPatterns[0], tpl)
		if tpl.DataStream != nil {
			ensureDatastream(ctx, tpl.IndexPatterns[0])
		}
	}
}

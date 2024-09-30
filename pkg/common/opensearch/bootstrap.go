package opensearch

import (
	"context"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func bootstrap(ctx context.Context) {
	for _, tpl := range api.INDEXES {
		ensureIndex(ctx, tpl.IndexPatterns[0], tpl)
		if tpl.DataStream != nil {
			ensureDatastream(ctx, tpl.IndexPatterns[0])
		}
	}
}

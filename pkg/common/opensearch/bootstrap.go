package opensearch

import (
	"context"

	"github.com/japannext/snooze/pkg/models"
)

func bootstrap(ctx context.Context) {
	for name, tpl := range models.OpensearchIndexTemplates {
		ensureIndexTemplate(ctx, name, tpl)

		if tpl.DataStream != nil {
			ensureDatastream(ctx, name)
		}
	}

	for name, idx := range models.OpensearchIndices {
		ensureIndice(ctx, name, idx)
	}
}

package opensearch

import (
	"context"

    dsl "github.com/mottaquikarim/esquerydsl"
    "github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const NOTIFICATION_INDEX = "v2-notifications"

func SearchNotifications(ctx context.Context, text string, timerange *api.TimeRange, pagination *api.Pagination) (*api.ListOf[*api.Notification], error) {
    var doc = &dsl.QueryDoc{}
    var params = &opensearchapi.SearchParams{}

	if pagination.OrderBy == "" {
		pagination.OrderBy = "timestampMillis"
	}

    addTimeRange(doc, timerange)
    addPagination(doc, params, pagination)
    addSearch(doc, text)

	return search[*api.Notification](ctx, NOTIFICATION_INDEX, params, doc)
}

func StoreNotification(ctx context.Context, item *api.Notification) (string, error) {
	return store(ctx, NOTIFICATION_INDEX, item)
}

func bootstrapNotifications(ctx context.Context) {
	numberOfShards := 1
	numberOfReplicas := 2
	settings := IndexSettings{numberOfShards, numberOfReplicas}
	tpl := IndexTemplate{
		IndexPattern: []string{NOTIFICATION_INDEX},
		DataStream: &DataStream{TimestampField{"timestampMillis"}},
		Template: Indice{
			Settings: settings,
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestampMillis": {Type: "date", Format: "epoch_millis"},
					"destination.kind": {Type: "keyword"},
					"destination.name": {Type: "keyword"},
				},
			},
		},
	}
	ensureIndex(ctx, NOTIFICATION_INDEX, tpl)
	ensureDatastream(ctx, NOTIFICATION_INDEX)
}

func init() {
	bootstraps = append(bootstraps, bootstrapNotifications)
}

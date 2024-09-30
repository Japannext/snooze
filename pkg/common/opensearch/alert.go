package opensearch

import (
    "context"

    dsl "github.com/mottaquikarim/esquerydsl"
    "github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const ALERT_INDEX = "v2-alerts"

func SearchAlerts(ctx context.Context, pagination *api.Pagination, history bool) (*api.ListOf[*api.Alert], error) {
	var doc = &dsl.QueryDoc{}
	var params = &opensearchapi.SearchParams{}

	if pagination.OrderBy == "" {
		pagination.OrderBy = "startsAt"
	}

	if !history {
		doc.And = []dsl.QueryItem{{Type: dsl.Term, Field: "endsAt", Value: 0}}
	}

	addPagination(doc, params, pagination)

	return search[*api.Alert](ctx, ALERT_INDEX, params, doc)
}

func UpdateAlert(ctx context.Context, item *api.Alert) error {
	// TODO
	return nil
}

func StoreAlert(ctx context.Context, item *api.Alert) (string, error) {
	return store(ctx, ALERT_INDEX, item)
}

func bootstrapAlerts (ctx context.Context) {
	numberOfShards := 1
	numberOfReplicas := 2
	settings := IndexSettings{numberOfShards, numberOfReplicas}
	alertTemplate := IndexTemplate{
		IndexPattern: []string{ALERT_INDEX},
		// DataStream: DataStream{TimestampField{"timestampMillis"}},
		Template: Indice{
			Settings: settings,
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"hash": {Type: "keyword"},
					"startsAt": {Type: "date", Format: "epoch_millis"},
					"endsAt": {Type: "date", Format: "epoch_millis"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"labels": {Type: "object"},
				},
			},
		},
	}
	ensureIndex(ctx, ALERT_INDEX, alertTemplate)
}

func init() {
	bootstraps = append(bootstraps, bootstrapAlerts)
}

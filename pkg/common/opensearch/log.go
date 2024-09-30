package opensearch

import (
	"context"

    dsl "github.com/mottaquikarim/esquerydsl"
    "github.com/opensearch-project/opensearch-go/v4/opensearchapi"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const LOG_INDEX = "v2-logs"

func SearchLogs(ctx context.Context, text string, timerange *api.TimeRange, pagination *api.Pagination) (*api.ListOf[*api.Log], error) {
    var doc = &dsl.QueryDoc{}
    var params = &opensearchapi.SearchParams{}

	if pagination.OrderBy == "" {
		pagination.OrderBy = "timestampMillis"
	}

    addTimeRange(doc, timerange)
    addPagination(doc, params, pagination)
    addSearch(doc, text)

	return search[*api.Log](ctx, LOG_INDEX, params, doc)
}

func StoreLog(ctx context.Context, item *api.Log) (string, error) {
	return store(ctx, LOG_INDEX, item)
}

func bootstrapLogs(ctx context.Context) {
	numberOfShards := 1
	numberOfReplicas := 2
	settings := IndexSettings{numberOfShards, numberOfReplicas}
	tpl := IndexTemplate{
		IndexPattern: []string{LOG_INDEX},
		DataStream: &DataStream{TimestampField{"timestampMillis"}},
		Template: Indice{
			Settings: settings,
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"timestampMillis": {Type: "date", Format: "epoch_millis"},
					"source.kind": {Type: "keyword"},
					"source.name": {Type: "keyword"},
					"identity": {Type: "object"},
					"group.hash":   {Type: "keyword"},
					"group.labels": {Type: "object"},
					"profile": {Type: "keyword"},
					"pattern": {Type: "keyword"},
					"labels":      {Type: "object"},
					"message":        {Type: "text"},
					"mute.skipNotification": {Type: "boolean"},
					"mute.skipStorage": {Type: "boolean"},
					"mute.reason": {Type: "keyword"},
				},
			},
		},
	}
	ensureIndex(ctx, LOG_INDEX, tpl)
	ensureDatastream(ctx, LOG_INDEX)
}

func init() {
	bootstraps = append(bootstraps, bootstrapLogs)
}

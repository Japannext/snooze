package opensearch

import (
	"context"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const LOG_INDEX = "v2-logs"

func SearchLogs(ctx context.Context, text string, timerange api.TimeRange, pagination api.Pagination) (*api.ListOf[*api.Log], error) {
	return search[*api.Log](ctx, LOG_INDEX, text, timerange, pagination)
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

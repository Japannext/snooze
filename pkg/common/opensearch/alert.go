package opensearch

import (
    "context"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

const ALERT_INDEX = "v2-alerts"

func SearchAlerts(ctx context.Context, text string, timerange api.TimeRange, pagination api.Pagination) (*api.ListOf[*api.Alert], error) {
	return search[*api.Alert](ctx, ALERT_INDEX, text, timerange, pagination)
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

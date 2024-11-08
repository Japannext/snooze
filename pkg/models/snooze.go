package models

const SNOOZE_INDEX = "v2-snoozes"

type Snooze struct {
	Base
	// Groups that will be snoozed
	Groups []Group `json:"groups"`
	Reason string `json:"reason"`
	StartAt Time `json:"startAt"`
	ExpireAt Time `json:"expireAt"`
}

func init() {
	index := IndexTemplate{
		Version: 1,
		IndexPatterns: []string{SNOOZE_INDEX},
		DataStream: map[string]map[string]string{"timestamp_field": {"name": "startAt"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"groups.name": {Type: "keyword"},
					"groups.hash": {Type: "keyword"},
					"reason": {Type: "text"},
					"startAt": {Type: "date", Format: "epoch_millis"},
					"expireAt": {Type: "date", Format: "epoch_millis"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}

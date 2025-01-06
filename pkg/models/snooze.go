package models

const SNOOZE_INDEX = "v2-snoozes"

type Snooze struct {
	Base
	// Groups that will be snoozed
	Groups    []Group       `json:"groups"`
	Reason    string        `json:"reason"`
	Tags      []Tag         `json:"tags"`
	StartsAt  Time          `json:"startAt"`
	EndsAt    Time          `json:"endsAt"`
	Cancelled *SnoozeCancel `json:"cancelled,omitempty"`
	Username  string        `json:"username"`
}

type SnoozeCancel struct {
	// The user that cancelled the snooze
	By string `json:"by"`
	// The time the snooze was cancelled
	At Time `json:"cancelledAt"`
	// The reason why the snooze was cancelled
	Reason string `json:"reason"`
}

func init() {
	index := IndexTemplate{
		Version:       3,
		IndexPatterns: []string{SNOOZE_INDEX},
		DataStream:    map[string]map[string]string{"timestamp_field": {"name": "startAt"}},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"groups.name":      {Type: "keyword"},
					"groups.hash":      {Type: "keyword"},
					"reason":           {Type: "text"},
					"startsAt":         {Type: "date", Format: "epoch_millis"},
					"endsAt":           {Type: "date", Format: "epoch_millis"},
					"tags.name":        {Type: "keyword"},
					"tags.color":       {Type: "keyword"},
					"tags.description": {Type: "text"},
					"cancelled.at":     {Type: "date", Format: "epoch_millis"},
					"cancelled.by":     {Type: "keyword"},
					"cancelled.reason": {Type: "text"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}

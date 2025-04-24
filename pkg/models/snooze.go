package models

const SnoozeIndex = "v2-snoozes"

var SnoozeIndices = IndexTemplate{
	Version:       3,
	IndexPatterns: []string{SnoozeIndex},
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

// A snooze entry that will be stored in Opensearch in order to search
// easily
type SnoozeEntry struct {
	Base
	// Groups that will be snoozed
	Groups    []Group       `json:"groups"`
	Reason    string        `json:"reason"`
	Tags      []Tag         `json:"tags"`
	StartsAt  Time          `json:"startAt"`
	EndsAt    Time          `json:"endsAt"`
	Cancelled *SnoozeCancel `json:"cancelled,omitempty"`
	Username  string        `json:"username"`
	If        string        `json:"if"`
}

// A snooze entry placed in redis for fast lookup of what
// snooze entry matches a log.
type SnoozeLookup struct {
	// ID of the associated Opensearch object (Snooze entry)
	OpensearchID string `json:"opensearchID"`
	// An additional condition of the snooze entry
	If string `json:"if"`
	// When the snooze entry should start. Only matter for edge cases.
	StartsAt Time `json:"startsAt"`
	// When the snooze entry should end. Only matter for edge cases.
	EndsAt Time `json:"endsAt"`
	// Number of times logs hit this snooze entry
	Hits int `json:"hits"`
}

type SnoozeCancel struct {
	// The user that cancelled the snooze
	By string `json:"by"`
	// The time the snooze was cancelled
	At Time `json:"cancelledAt"`
	// The reason why the snooze was cancelled
	Reason string `json:"reason"`
}

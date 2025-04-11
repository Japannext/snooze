package models

const RatelimitHistoryIndex = "v2-ratelimit-history"

var RatelimitHistoryIndexTemplate = IndexTemplate{
	Version:       0,
	IndexPatterns: []string{RatelimitHistoryIndex},
	DataStream:    map[string]map[string]string{"timestamp_field": {"name": "startsAt"}},
	Template: Indice{
		Settings: IndexSettings{1, 2},
		Mappings: IndexMapping{
			Properties: map[string]MappingProps{
				"startsAt": {Type: "date", Format: "epoch_millis"},
				"endsAt":   {Type: "date", Format: "epoch_millis"},
				"hash":     {Type: "keyword"},
				"rule":     {Type: "keyword"},
			},
		},
	},
}

// A rate limit status. Stored in Redis, and
type RatelimitStatus struct {
	// Whether the rate limiting is active or over
	Active bool `json:"active"`
	// Time when the entries were rate-limited or allowed back
	StartsAt Time `json:"startsAt"`
	// Time when the ratelimit should be retried
	RetryAt uint64 `json:"retryAt"`
	// Time when the concerned entries stopped being rate-limited
	EndsAt uint64 `json:"endsAt"`
	// The name of the ratelimit rule that was applied
	Rule string `json:"name"`
	// The hash of the fields. Can uniquely identify a ratelimited identity
	Hash string `json:"hash"`
}

// An entry in the rate-limit history (persistent storage)
type RatelimitRecord struct {
	Base

	// Whether the rate limiting is active or over
	Active   bool `json:"active"`
	StartsAt uint64
	EndsAt   uint64
	Rule     string
	Hash     string
}

type RatelimitStatusResults struct {
	Items []RatelimitStatus `json:"items"`
	Total int               `json:"total"`
	More  bool              `json:"more"`
}

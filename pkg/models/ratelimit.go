package models

const RatelimitIndex = "v2-ratelimits"

// A rate limit object. Represent the start or end of rate limiting.
type Ratelimit struct {
	ID string `json:"id"`
	// Time when the entries were rate-limited or allowed back
	StartsAt uint64 `yaml:"startsAt"`
	// Time when the concerned entries stopped being rate-limited
	EndsAt uint64 `yaml:"endsAt"`
	// The name of the ratelimit rule that was applied
	Rule string `json:"name"`
	// The fields that are concerned by the rate limiting
	Fields map[string]string `yaml:"fields"`
	// The hash of the fields. Can uniquely identify a ratelimited identity
	Hash string `yaml:"hash"`
}

type RatelimitResults struct {
	Items []Ratelimit `json:"items"`
	Total int         `json:"total"`
	More  bool        `json:"more"`
}

func init() {
	OpensearchIndexTemplates[RatelimitIndex] = IndexTemplate{
		Version:       0,
		IndexPatterns: []string{RatelimitIndex},
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
}

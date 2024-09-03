package v2

// A rate limit object. Represent the start or end of rate limiting.
type Ratelimit struct {
	ID string `json:"id"`
	// Valid values: "start"|"end"
	Kind string `json:"kind"`
	// Time when the entries were ratelimited or allowed back
	TimestampMillis uint64 `yaml:"timestampMillis"`
	// The name of the ratelimit rule that was applied
	Rule string `json:"name"`
	// The fields that are concerned by the rate limiting
	Fields map[string]string `yaml:"fields"`
	// The hash of the fields. Can uniquely identify a ratelimited identity
	Hash string `yaml:"hash"`
}

type RatelimitResults struct {
	Items []Ratelimit `json:"items"`
	Total int `json:"total"`
	More bool `json:"more"`
}

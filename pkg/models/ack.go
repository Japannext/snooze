package models

const AckIndex = "v2-ack"

// An acknowledgement given to a log or
// group of logs.
type Ack struct {
	Base
	Time     Time     `json:"timestamp"`
	Username string   `json:"username"`
	Reason   string   `json:"reason"`
	LogIDs   []string `json:"logIDs"`
}

func init() {
	OpensearchIndices[AckIndex] = Indice{
		Settings: IndexSettings{1, 2},
		Mappings: IndexMapping{
			Properties: map[string]MappingProps{
				"time":     {Type: "date", Format: "epoch_millis"},
				"username": {Type: "keyword"},
				"reason":   {Type: "text"},
				"logIDs":   {Type: "keyword"},
			},
		},
	}
}

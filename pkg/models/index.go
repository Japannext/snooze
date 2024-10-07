package models

var INDEXES []IndexTemplate

type IndexTemplate struct {
	Version int `json:"version"`
    IndexPatterns []string `json:"index_patterns"`
    DataStream map[string]map[string]string `json:"data_stream,omitempty"`
    Template     Indice   `json:"template"`
}

type Indice struct {
    Settings IndexSettings  `json:"settings"`
    Mappings IndexMapping `json:"mappings"`
}

type IndexSettings struct {
    NumberOfShards   int `json:"number_of_shards"`
    NumberOfReplicas int `json:"number_of_replicas"`
}

type MappingProps struct {
    Type   string                  `json:"type,omitempty"`
    Format string                  `json:"format,omitempty"`
    Fields map[string]MappingProps `json:"fields,omitempty"`
}

type IndexMapping struct {
    Properties map[string]MappingProps `json:"properties"`
}

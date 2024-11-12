package models

const GROUP_INDEX = "v2-groups"

type Group struct {
	Base

	Name string `json:"name"`
	// Hash value of the key-value
	Hash string `json:"hash"`
	// Human readable information about the group
	Labels map[string]string `json:"labels,omitempty"`
}

func init() {
	index := IndexTemplate{
		Version: 1,
		IndexPatterns: []string{GROUP_INDEX},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"name": {Type: "keyword"},
					"labels": {Type: "object"},
					"hash": {Type: "keyword"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}

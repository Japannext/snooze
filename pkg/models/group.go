package models

const GROUP_INDEX = "v2-groups"

type Group struct {
	Base

	Name string `json:"name"`
	// Hash value of the key-value
	Hash string `json:"hash"`
	// Human readable information about the group
	Labels map[string]string `json:"labels,omitempty"`

	// Last time a new hit of this group was inserted.
	// This is mainly used to curate the index
	LastInsert Time `json:"lastInsert"`
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
					"lastInsert": {Type: "date", Format: "epoch_millis"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}

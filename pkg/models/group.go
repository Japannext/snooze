package models

const GroupIndex = "v2-groups"

var GroupIndice = Indice{
	Settings: IndexSettings{1, 2},
	Mappings: IndexMapping{
		Properties: map[string]MappingProps{
			"name":   {Type: "keyword"},
			"labels": {Type: "object"},
			"hash":   {Type: "keyword"},
		},
	},
}

type Group struct {
	Base

	Name string `json:"name"`
	// Hash value of the key-value
	Hash string `json:"hash"`
	// Human readable information about the group
	Labels map[string]string `json:"labels,omitempty"`
}

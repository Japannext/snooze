package models

const TagIndex = "v2-tags"

var TagIndice = Indice{
	Settings: IndexSettings{1, 2},
	Mappings: IndexMapping{
		Properties: map[string]MappingProps{
			"name":        {Type: "keyword"},
			"description": {Type: "text"},
			"color":       {Type: "keyword"},
		},
	},
}

type Tag struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

package models

const TAG_INDEX = "v2-tags"

type Tag struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Color       string `json:"color,omitempty"`
}

func init() {
	index := IndexTemplate{
		Version:       3,
		IndexPatterns: []string{TAG_INDEX},
		Template: Indice{
			Settings: IndexSettings{1, 2},
			Mappings: IndexMapping{
				Properties: map[string]MappingProps{
					"name":        {Type: "keyword"},
					"description": {Type: "text"},
					"color":       {Type: "keyword"},
				},
			},
		},
	}
	INDEXES = append(INDEXES, index)
}

package dsl

/*
type Query struct {
	Index       string `json:"index"`
	Size        int `json:"size"`
	From        int `json:from"`
	Sort        []map[string]string `json:"sort"`
	SearchAfter []interface{}
	And         []QueryItem
	Not         []QueryItem
	Or          []QueryItem
	Filter      []QueryItem
	PageSize    int
}
*/

type QueryRequest struct {
	Query Query `json:"query"`
}

type Query struct {
	Term map[string]string `json:"term,omitempty"`

}

func NewTermQuery(field, value string) *Query {
	return &Query{Term: map[string]string{field: value}}
}

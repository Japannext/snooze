package dsl

type QueryRequest struct {
	Query Query `json:"query"`
	Sort []map[string]Sort `json:"sort"`
}

type Query struct {
	Bool *Bool `json:"bool,omitempty"`
}

type Sort struct {
	Order string `json:"order"`
}

type Bool struct {
	And []QueryItem `json:"must,omitempty"`
	Not []QueryItem `json:"must_not,omitempty"`
	Or []QueryItem `json:"should,omitempty"`
}

type QueryItem struct {
	Term map[string]string `json:"term,omitempty"`
	Range map[string]Range `json:"range,omitempty"`
	Match map[string]Match `json:"match,omitempty"`
	QueryString *QueryString `json:"query_string,omitempty"`
}

type QueryString struct {
	Query string `json:"query,omitempty"`
}

type Match struct {
	Query string `json:"query,omitempty"`
	Operator string `json:"operator,omitempty"`
	MinimumShouldMatch int `json:"minimum_should_match,omitempty"`
	Analyzer string `json:"analyzer,omitempty"`
}

type Range struct {
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func (query *Query) WithTerm(key, value string) *Query {
	item := QueryItem{Term: map[string]string{key: value}}
	query.Bool.And = append(query.Bool.And, item)
	return query
}

package dsl

type QueryReq struct {
	Query `json:"query"`
	Sort []map[string]Sort `json:"sort,omitempty"`
	Script *Script `json:"script,omitempty"`
}

/*
type QueryRequest struct {
	Query Query `json:"query"`
	Sort []map[string]Sort `json:"sort,omitempty"`
	Script *Script `json:"script,omitempty"`
}
*/

type Query struct {
	Bool *Bool `json:"bool,omitempty"`
}

type Sort struct {
	Order string `json:"order"`
}

type Bool struct {
	Must []QueryItem `json:"must,omitempty"`
	MustNot []QueryItem `json:"must_not,omitempty"`
	Should []QueryItem `json:"should,omitempty"`
	MinimumShouldMatch int `json:"minimum_should_match,omitempty"`
}

type Script struct {
	Source string `json:"source"`
	Lang string `json:"lang"`
	Params map[string]interface{} `json:"params"`
}

type QueryItem struct {
	Term map[string]interface{} `json:"term,omitempty"`
	Terms map[string][]string `json:"terms,omitempty"`
	Range map[string]Range `json:"range,omitempty"`
	Match map[string]Match `json:"match,omitempty"`
	Exists *Exists `json:"exists,omitempty"`
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

type Exists struct {
	Field string `json:"field"`
}

type Range struct {
	Gte *uint64 `json:"gte,omitempty"`
	Lte *uint64 `json:"lte,omitempty"`
}

func (query *Query) WithTerm(key string, value interface{}) *Query {
	item := QueryItem{Term: map[string]interface{}{key: value}}
	return query.And(item)
}

func (query *Query) WithTerms(key string, values []string) *Query {
	item := QueryItem{Terms: map[string][]string{key: values}}
	return query.And(item)
}

func (query *Query) And(item QueryItem) *Query {
	if query.Bool == nil {
		query.Bool = &Bool{}
	}
	query.Bool.Must = append(query.Bool.Must, item)
	return query
}

func (query *Query) Or(item QueryItem) *Query {
	if query.Bool == nil {
		query.Bool = &Bool{}
	}
	query.Bool.Should = append(query.Bool.Should, item)
	return query
}

func (query *Query) WithExists(field string) *Query {
	item := QueryItem{Exists: &Exists{Field: field}}
	query.And(item)
	return query
}

func (query *Query) WithQueryString(text string) *Query {
	item := QueryItem{QueryString: &QueryString{Query: text}}
	query.And(item)
	return query
}

func (req *QueryReq) WithSort(field string, ascending bool) *QueryReq {
	var order Sort
	if ascending {
		order = Sort{Order: "asc"}
	} else {
		order = Sort{Order: "desc"}
	}
	req.Sort = append(req.Sort, map[string]Sort{field: order})
	return req
}

func (req *QueryReq) WithPainlessScript(script string, params map[string]interface{}) *QueryReq {
	req.Script = &Script{
		Source: script,
		Lang: "painless",
		Params: params,
	}
	return req
}

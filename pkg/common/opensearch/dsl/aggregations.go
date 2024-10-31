package dsl

type AggregationRequest struct {
	Aggs map[string]AggItem `json:"aggs"`
}

type AggItem struct {
	Terms *AggTerms `json:"terms"`
	Filter *AggFilter `json:"-,squash"`
}

type AggTerms struct {
	Field string `json:"field"`
}

type AggFilter struct {
	Filter *Query `json:"filter"`
	Aggs map[string]AggItem `json:"aggs"`
}

type AggregationResponse = map[string]AggRespItem

type AggRespItem struct {
	Buckets []Bucket `json:"buckets"`
}

type Bucket struct {
	Key string `json:"key"`
	DocCount int `json:"doc_count"`
}

func NewAggregationQuery(field string, filter *Query) *AggregationRequest {
	return &AggregationRequest{
		Aggs: map[string]AggItem{
			"filter-with": {
				Filter: &AggFilter{
					Filter: filter,
					Aggs: map[string]AggItem{
						"by-field": {
							Terms: &AggTerms{
								Field: field,
							},
						},
					},
				},
			},
		},
	}
}

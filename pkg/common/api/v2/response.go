package v2

type HasID interface {
	GetID() string
	SetID(string)
}

type ListOf[T HasID] struct {
	Items []*T `json:"items"`
	Total int `json:"total"`
	More bool `json:"more"`
}

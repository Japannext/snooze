package models

import (
	"github.com/google/uuid"
)

type Base struct {
	ID string `json:"_id,omitempty"`
}

func (item *Base) GetID() string { return item.ID }
func (item *Base) SetID(id string) { item.ID = id }

func (item *Base) NewID() {
	item.ID = uuid.NewString()
}

type HasContext interface {
	Context() map[string]interface{}
}

type HasID interface {
	GetID() string
	SetID(string)
}

type ListOf[T HasID] struct {
	Items []T `json:"items"`
	Total int `json:"total"`
	More bool `json:"more"`
}

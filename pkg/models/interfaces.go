package models

type Object interface {
	GetID() string
	SetID(id string)
	Context() map[string]interface{}
}

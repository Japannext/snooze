package models

type Base struct {
	ID string `json:"_id,omitempty"`
}

func (item Base) GetID() string { return item.ID }
func (item Base) SetID(id string) { item.ID = id }

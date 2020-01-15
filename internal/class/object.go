package class

import "github.com/google/uuid"

type Object interface {
	UUID() *string
	SetNewId()
}

type SerialObject struct {
	Id string `json:"id"` // UUID
}

func (o *SerialObject) UUID() *string {
	return &o.Id
}

func (o *SerialObject) SetNewId() {
	s, _ := uuid.NewUUID()
	o.Id = s.String()
}

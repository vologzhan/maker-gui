package request

import "github.com/google/uuid"

type EntityCreate struct {
	ServiceId uuid.UUID `param:"service-id"`
	Id        uuid.UUID `json:"id"`
	NameDb    string    `json:"name_db"`
}

type EntityDelete struct {
	Id uuid.UUID `param:"id"`
}

type EntityUpdate struct {
	Id     uuid.UUID `param:"id"`
	NameDb string    `json:"name_db"`
}

type EntityList struct {
	ServiceId uuid.UUID `param:"service-id"`
}

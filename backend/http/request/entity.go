package request

import "github.com/google/uuid"

type EntityCreate struct {
	ServiceId  uuid.UUID `param:"service-id"`
	Id         uuid.UUID `json:"id"`
	NameDb     string    `json:"name_db"`
	NamePlural string    `json:"name_plural"`
}

type EntityDelete struct {
	Id uuid.UUID `param:"id"`
}

type EntityUpdate struct {
	Id         uuid.UUID `param:"id"`
	NameDb     string    `json:"name_db"`
	NamePlural string    `json:"name_plural"`
}

type EntityList struct {
	ServiceId uuid.UUID `param:"service-id"`
}

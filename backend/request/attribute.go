package request

import "github.com/google/uuid"

type AttributeCreate struct {
	EntityId   uuid.UUID `param:"entity-id"`
	Id         uuid.UUID `json:"id"`
	NameDb     string    `json:"name_db"`
	TypeDb     string    `json:"type_db"`
	Default    string    `json:"default"`
	FkTable    string    `json:"fk_table"`
	FkType     string    `json:"fk_type"`
	Nullable   bool      `json:"nullable"`
	PrimaryKey bool      `json:"primary_key"`
}

type AttributeDelete struct {
	Id uuid.UUID `param:"id"`
}

type AttributeUpdate struct {
	Id         uuid.UUID `param:"id"`
	NameDb     string    `json:"name_db"`
	TypeDb     string    `json:"type_db"`
	Default    string    `json:"default"`
	FkTable    string    `json:"fk_table"`
	FkType     string    `json:"fk_type"`
	Nullable   bool      `json:"nullable"`
	PrimaryKey bool      `json:"primary_key"`
}

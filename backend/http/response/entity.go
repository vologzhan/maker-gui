package response

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker-gui/backend/entity"
)

type EntityList struct {
	Items []Entity `json:"items"`
}

type Entity struct {
	Id         uuid.UUID   `json:"id"`
	NameDb     string      `json:"name_db"`
	Attributes []Attribute `json:"attributes"`
}

func NewEntityList(entities []*entity.Entity) *EntityList {
	items := make([]Entity, len(entities))
	for i, e := range entities {
		items[i] = newEntity(e)
	}

	return &EntityList{items}
}

func newEntity(e *entity.Entity) Entity {
	return Entity{
		e.Id(),
		e.NameDb(),
		newAttributeSlice(e.MustAttributes()),
	}
}

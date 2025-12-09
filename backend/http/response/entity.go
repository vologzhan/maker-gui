package response

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker-common/strcase"
	"github.com/vologzhan/maker-gui/backend/maker/models"
)

type EntityList struct {
	Items []Entity `json:"items"`
}

type Entity struct {
	Id         uuid.UUID   `json:"id"`
	NameDb     string      `json:"name_db"`
	NamePlural string      `json:"name_plural"`
	Attributes []Attribute `json:"attributes"`
}

func NewEntityList(entities []*models.Entity) *EntityList {
	items := make([]Entity, len(entities))
	for i, entity := range entities {
		items[i] = newEntity(entity)
	}

	return &EntityList{items}
}

func newEntity(entity *models.Entity) Entity {
	return Entity{
		entity.Id(),
		entity.NameDb(),
		strcase.ToSnake(entity.NamePlural()),
		newAttributeSlice(entity.MustAttributes()),
	}
}

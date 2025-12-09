package response

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker-gui/backend/maker/models"
)

type Attribute struct {
	Id         uuid.UUID `json:"id"`
	NameDb     string    `json:"name_db"`
	TypeDb     string    `json:"type_db"`
	Default    string    `json:"default"`
	FkTable    string    `json:"fk_table"`
	FkType     string    `json:"fk_type"`
	Nullable   bool      `json:"nullable"`
	PrimaryKey bool      `json:"primary_key"`
}

func newAttributeSlice(attrs []*models.Attribute) []Attribute {
	items := make([]Attribute, len(attrs))
	for i, attr := range attrs {
		items[i] = newAttribute(attr)
	}

	return items
}

func newAttribute(attr *models.Attribute) Attribute {
	return Attribute{
		Id:         attr.Id(),
		NameDb:     attr.NameDb(),
		TypeDb:     attr.TypeDb(),
		Default:    attr.Default(),
		FkTable:    attr.FkTable(),
		FkType:     attr.FkType(),
		Nullable:   attr.Nullable(),
		PrimaryKey: attr.PrimaryKey(),
	}
}

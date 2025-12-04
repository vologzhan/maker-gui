package entity

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker"
	"github.com/vologzhan/maker-gui/backend/internal/slices"
	"github.com/vologzhan/maker-gui/backend/internal/typeconv"
)

const (
	attrNameDb   = "name_db"
	attrTypeDb   = "type_db"
	attrDefault  = "default"
	attrFkTable  = "fk_table"
	attrFkColumn = "fk_column"
	attrNullable = "nullable"
	attrPk       = "primary_key"
)

type Attribute struct {
	node   *maker.Node
	entity *Entity
}

func (a *Attribute) Id() uuid.UUID    { return a.node.Id() }
func (a *Attribute) NameDb() string   { return a.node.ValueString(attrNameDb) }
func (a *Attribute) TypeDb() string   { return a.node.ValueString(attrTypeDb) }
func (a *Attribute) Default() string  { return a.node.ValueString(attrDefault) }
func (a *Attribute) FkTable() string  { return a.node.ValueString(attrFkTable) }
func (a *Attribute) FkColumn() string { return a.node.ValueString(attrFkColumn) }
func (a *Attribute) Nullable() bool   { return a.node.ValueBool(attrNullable) }
func (a *Attribute) PrimaryKey() bool { return a.node.ValueBool(attrPk) }
func (a *Attribute) Entity() *Entity  { return a.entity }

func (a *Attribute) Delete() error {
	err := a.node.Delete()
	if err != nil {
		return err
	}

	slices.Delete(a.entity.attributes, a)

	return nil
}

func (a *Attribute) Update(nameDb, typeDb, def, fkTable, fkColumn string, nullable, pk bool) error {
	return a.node.SetValues(newAttributeValues(nameDb, typeDb, def, fkTable, fkColumn, nullable, pk))
}

func NewAttribute(entity *Entity, id uuid.UUID, nameDb, typeDb, def, fkTable, fkColumn string, nullable, pk bool) (*Attribute, error) {
	node, err := entity.node.CreateChild("attribute", id, newAttributeValues(nameDb, typeDb, def, fkTable, fkColumn, nullable, pk))
	if err != nil {
		return nil, err
	}

	return newAttribute(node, entity), nil
}

func newAttribute(node *maker.Node, entity *Entity) *Attribute {
	attr := &Attribute{
		node,
		entity,
	}
	entity.attributes = append(entity.attributes, attr)

	return attr
}

func newAttributeValues(nameDb, typeDb, def, fkTable, fkColumn string, nullable, pk bool) map[string]string {
	return map[string]string{
		"name":        nameDb,
		"name_db":     nameDb,
		"type_db":     typeDb,
		"type_go":     typeconv.DbToGo(typeDb),
		"default":     def,
		"fk_table":    fkTable,
		"fk_column":   fkColumn,
		"nullable":    typeconv.BoolToString(nullable),
		"primary_key": typeconv.BoolToString(pk),
	}
}

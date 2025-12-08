package entity

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker"
	"github.com/vologzhan/maker/helper/slices"
)

const entityNameDb = "name_db"

type Entity struct {
	node       *maker.Node
	service    *Service
	attributes []*Attribute
}

func (e *Entity) Id() uuid.UUID     { return e.node.Id() }
func (e *Entity) NameDb() string    { return e.node.ValueString(entityNameDb) }
func (e *Entity) Service() *Service { return e.service }

func (e *Entity) Attributes() ([]*Attribute, error) {
	if e.attributes == nil {
		nodes, err := e.node.Children("attribute")
		if err != nil {
			return nil, err
		}

		e.attributes = make([]*Attribute, 0, len(nodes))

		for _, node := range nodes {
			_ = newAttribute(node, e)
		}
	}

	return e.attributes, nil
}

func (e *Entity) MustAttributes() []*Attribute {
	if e.attributes == nil {
		panic("entity: Entity.MustAttributes: attributes are requested before reading")
	}

	return e.attributes
}

func (e *Entity) Delete() error {
	err := e.node.Delete()
	if err != nil {
		return err
	}

	e.service.entities = slices.Delete(e.service.entities, e)

	return nil
}

func (e *Entity) Update(nameDb string) error {
	return e.node.SetValues(newEntityValues(nameDb))
}

func NewEntity(service *Service, id uuid.UUID, nameDb string) (*Entity, error) {
	node, err := service.node.CreateChild("entity", id, newEntityValues(nameDb))
	if err != nil {
		return nil, err
	}

	return newEntity(node, service, []*Attribute{}), nil
}

func newEntity(node *maker.Node, service *Service, attributes []*Attribute) *Entity {
	entity := &Entity{
		node,
		service,
		attributes,
	}
	service.entities = append(service.entities, entity)

	return entity
}

func newEntityValues(nameDb string) map[string]string {
	return map[string]string{
		entityNameDb: nameDb,
		"name":       nameDb,
	}
}

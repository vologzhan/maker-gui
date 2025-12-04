package entity

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker"
)

const serviceName = "name"

type Service struct {
	node     *maker.Node
	root     *Root
	sql      *Sql
	entities []*Entity
}

func (s *Service) Id() uuid.UUID { return s.node.Id() }
func (s *Service) Name() string  { return s.node.ValueString(serviceName) }
func (s *Service) Flush() error  { return s.node.Flush() }

func (s *Service) Entities() ([]*Entity, error) {
	if s.entities == nil {
		nodes, err := s.node.Children("entity")
		if err != nil {
			return nil, err
		}

		s.entities = make([]*Entity, 0, len(nodes))

		for _, node := range nodes {
			_ = newEntity(node, s, nil)
		}
	}

	return s.entities, nil
}

func NewService(root *Root, id uuid.UUID, name string) (*Service, error) {
	node, err := root.node.CreateChild("service", id, map[string]string{
		serviceName: name,
	})
	if err != nil {
		return nil, err
	}

	return newService(node, root, []*Entity{}), nil
}

func newService(node *maker.Node, root *Root, entities []*Entity) *Service {
	service := &Service{
		node,
		root,
		nil,
		entities,
	}
	root.services = append(root.services, service)

	return service
}

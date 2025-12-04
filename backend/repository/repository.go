package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vologzhan/maker"
	"github.com/vologzhan/maker-gui/backend/entity"
)

type Repository struct {
	root       *entity.Root
	services   map[uuid.UUID]*entity.Service
	entities   map[uuid.UUID]*entity.Entity
	attributes map[uuid.UUID]*entity.Attribute
}

func New(mak *maker.Node) *Repository {
	return &Repository{
		entity.NewRoot(mak),
		make(map[uuid.UUID]*entity.Service),
		make(map[uuid.UUID]*entity.Entity),
		make(map[uuid.UUID]*entity.Attribute),
	}
}

func (r *Repository) Root() *entity.Root {
	return r.root
}

func (r *Repository) Service(id uuid.UUID) (*entity.Service, error) {
	s, ok := r.services[id]
	if !ok {
		return nil, fmt.Errorf("repository: service '%s' not found", id)
	}

	return s, nil
}

func (r *Repository) CreateService(s *entity.Service) {
	if _, ok := r.services[s.Id()]; !ok {
		r.services[s.Id()] = s
	}
}

func (r *Repository) Entity(id uuid.UUID) (*entity.Entity, error) {
	e, ok := r.entities[id]
	if !ok {
		return nil, fmt.Errorf("repository: entity '%s' not found", id)
	}

	return e, nil
}

func (r *Repository) CreateEntity(e *entity.Entity) {
	if _, ok := r.entities[e.Id()]; !ok {
		r.entities[e.Id()] = e
	}
}

func (r *Repository) DeleteEntity(e *entity.Entity) {
	delete(r.entities, e.Id())
}

func (r *Repository) Attribute(id uuid.UUID) (*entity.Attribute, error) {
	a, ok := r.attributes[id]
	if !ok {
		return nil, fmt.Errorf("repository: attribute '%s' not found", id)
	}

	return a, nil
}

func (r *Repository) CreateAttribute(a *entity.Attribute) {
	if _, ok := r.attributes[a.Id()]; !ok {
		r.attributes[a.Id()] = a
	}
}

func (r *Repository) DeleteAttribute(a *entity.Attribute) {
	delete(r.attributes, a.Id())
}

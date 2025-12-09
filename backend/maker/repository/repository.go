package repository

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/vologzhan/maker"
	models "github.com/vologzhan/maker-gui/backend/maker/models"
)

type Repository struct {
	root       *models.Root
	services   map[uuid.UUID]*models.Service
	entities   map[uuid.UUID]*models.Entity
	attributes map[uuid.UUID]*models.Attribute
}

func New(mak *maker.Node) *Repository {
	return &Repository{
		models.NewRoot(mak),
		make(map[uuid.UUID]*models.Service),
		make(map[uuid.UUID]*models.Entity),
		make(map[uuid.UUID]*models.Attribute),
	}
}

func (r *Repository) Root() *models.Root {
	return r.root
}

func (r *Repository) Service(id uuid.UUID) (*models.Service, error) {
	s, ok := r.services[id]
	if !ok {
		return nil, fmt.Errorf("repository: service '%s' not found", id)
	}

	return s, nil
}

func (r *Repository) CreateService(s *models.Service) {
	if _, ok := r.services[s.Id()]; !ok {
		r.services[s.Id()] = s
	}
}

func (r *Repository) Entity(id uuid.UUID) (*models.Entity, error) {
	e, ok := r.entities[id]
	if !ok {
		return nil, fmt.Errorf("repository: entity '%s' not found", id)
	}

	return e, nil
}

func (r *Repository) CreateEntity(e *models.Entity) {
	if _, ok := r.entities[e.Id()]; !ok {
		r.entities[e.Id()] = e
	}
}

func (r *Repository) DeleteEntity(e *models.Entity) {
	delete(r.entities, e.Id())
}

func (r *Repository) Attribute(id uuid.UUID) (*models.Attribute, error) {
	a, ok := r.attributes[id]
	if !ok {
		return nil, fmt.Errorf("repository: attribute '%s' not found", id)
	}

	return a, nil
}

func (r *Repository) CreateAttribute(a *models.Attribute) {
	if _, ok := r.attributes[a.Id()]; !ok {
		r.attributes[a.Id()] = a
	}
}

func (r *Repository) DeleteAttribute(a *models.Attribute) {
	delete(r.attributes, a.Id())
}

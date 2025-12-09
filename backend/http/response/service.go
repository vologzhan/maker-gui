package response

import (
	"github.com/google/uuid"
	"github.com/vologzhan/maker-common/strcase"
	"github.com/vologzhan/maker-gui/backend/entity"
)

type ServiceList struct {
	Items []Service `json:"items"`
}

type Service struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func NewServiceList(services []*entity.Service) *ServiceList {
	items := make([]Service, len(services))
	for i, service := range services {
		items[i] = newService(service)
	}

	return &ServiceList{items}
}

func newService(service *entity.Service) Service {
	return Service{
		service.Id(),
		strcase.ToKebab(service.Name()),
	}
}

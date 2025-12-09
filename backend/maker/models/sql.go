package models

import "github.com/vologzhan/maker"

type Sql struct {
	node    *maker.Node
	service *Service
}

func NewSql(node *maker.Node, service *Service) *Sql {
	return &Sql{
		node,
		service,
	}
}

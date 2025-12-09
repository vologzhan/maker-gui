package models

import (
	"github.com/vologzhan/maker"
)

type Root struct {
	node     *maker.Node
	services []*Service
}

func (r *Root) Services() ([]*Service, error) {
	if r.services == nil {
		nodes, err := r.node.Children("service")
		if err != nil {
			return nil, err
		}

		r.services = make([]*Service, 0, len(nodes))

		for _, node := range nodes {
			_ = newService(node, r, nil)
		}
	}

	return r.services, nil
}

func NewRoot(node *maker.Node) *Root {
	return &Root{node, nil}
}

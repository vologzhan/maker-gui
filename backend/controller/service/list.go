package service

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/repository"
	"github.com/vologzhan/maker-gui/backend/response"
	"net/http"
)

type List struct {
	rep *repository.Repository
}

func (c *List) Handle(ctx echo.Context) error {
	res, err := c.handle()
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (c *List) handle() (*response.ServiceList, error) {
	root := c.rep.Root()

	services, err := root.Services()
	if err != nil {
		return nil, err
	}

	for _, service := range services {
		c.rep.CreateService(service)
	}

	return response.NewServiceList(services), nil
}

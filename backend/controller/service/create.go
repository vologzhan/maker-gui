package service

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/entity"
	"github.com/vologzhan/maker-gui/backend/repository"
	"github.com/vologzhan/maker-gui/backend/request"
	"github.com/vologzhan/maker-gui/backend/response"
	"net/http"
)

type Create struct {
	rep *repository.Repository
}

func (c *Create) Handle(ctx echo.Context) error {
	var req request.ServiceCreate
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Create) handle(req request.ServiceCreate) error {
	root := c.rep.Root()

	service, err := entity.NewService(root, req.Id, req.Name)
	if err != nil {
		return err
	}

	c.rep.CreateService(service)

	return service.Flush()
}

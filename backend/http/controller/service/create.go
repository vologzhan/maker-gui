package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/maker/models"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
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

	service, err := models.NewService(root, req.Id, req.Name)
	if err != nil {
		return err
	}

	c.rep.CreateService(service)

	return service.Flush()
}

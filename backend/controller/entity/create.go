package entity

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/entity"
	"github.com/vologzhan/maker-gui/backend/repository"
	"github.com/vologzhan/maker-gui/backend/request"
	"github.com/vologzhan/maker-gui/backend/response"
	"net/http"
)

type Create struct {
	repository *repository.Repository
}

func (c *Create) Handle(ctx echo.Context) error {
	var req request.EntityCreate
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Create) handle(req request.EntityCreate) error {
	service, err := c.repository.Service(req.ServiceId)
	if err != nil {
		return err
	}

	e, err := entity.NewEntity(service, req.Id, req.NameDb)
	if err != nil {
		return err
	}

	c.repository.CreateEntity(e)

	// todo sql

	return service.Flush()
}

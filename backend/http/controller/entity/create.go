package entity

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/maker/models"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
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

	entity, err := models.NewEntity(service, req.Id, req.NameDb, req.NamePlural)
	if err != nil {
		return err
	}

	c.repository.CreateEntity(entity)

	// todo sql

	return service.Flush()
}

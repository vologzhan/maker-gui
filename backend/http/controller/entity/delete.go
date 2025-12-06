package entity

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/repository"
	"net/http"
)

type Delete struct {
	repository *repository.Repository
}

func (c *Delete) Handle(ctx echo.Context) error {
	var req request.EntityDelete
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Delete) handle(req request.EntityDelete) error {
	entity, err := c.repository.Entity(req.Id)
	if err != nil {
		return err
	}

	err = entity.Delete()
	if err != nil {
		return err
	}

	c.repository.DeleteEntity(entity)

	service := entity.Service()

	// todo sql

	return service.Flush()
}

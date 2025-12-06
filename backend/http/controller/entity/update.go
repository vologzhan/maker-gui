package entity

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/repository"
	"net/http"
)

type Update struct {
	repository *repository.Repository
}

func (c *Update) Handle(ctx echo.Context) error {
	var req request.EntityUpdate
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Update) handle(req request.EntityUpdate) error {
	entity, err := c.repository.Entity(req.Id)
	if err != nil {
		return err
	}

	err = entity.Update(req.NameDb)
	if err != nil {
		return err
	}

	service := entity.Service()

	// todo sql

	return service.Flush()
}

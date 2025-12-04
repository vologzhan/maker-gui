package attribute

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/repository"
	"github.com/vologzhan/maker-gui/backend/request"
	"github.com/vologzhan/maker-gui/backend/response"
	"net/http"
)

type Delete struct {
	repository *repository.Repository
}

func (c *Delete) Handle(ctx echo.Context) error {
	var req request.AttributeDelete
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Delete) handle(req request.AttributeDelete) error {
	attr, err := c.repository.Attribute(req.Id)
	if err != nil {
		return err
	}

	err = attr.Delete()
	if err != nil {
		return err
	}

	c.repository.DeleteAttribute(attr)

	service := attr.Entity().Service()

	// todo sql

	return service.Flush()
}

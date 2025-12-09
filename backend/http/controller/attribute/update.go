package attribute

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
	"net/http"
)

type Update struct {
	repository *repository.Repository
}

func (c *Update) Handle(ctx echo.Context) error {
	var req request.AttributeUpdate
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Update) handle(req request.AttributeUpdate) error {
	attr, err := c.repository.Attribute(req.Id)
	if err != nil {
		return err
	}

	err = attr.Update(req.NameDb, req.TypeDb, req.Default, req.FkTable, req.FkType, req.Nullable, req.PrimaryKey)
	if err != nil {
		return err
	}

	service := attr.Entity().Service()

	// todo sql

	return service.Flush()
}

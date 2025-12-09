package attribute

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/maker/models"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
	"net/http"
)

type Create struct {
	repository *repository.Repository
}

func (c *Create) Handle(ctx echo.Context) error {
	var req request.AttributeCreate
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewSuccess())
}

func (c *Create) handle(req request.AttributeCreate) error {
	entity, err := c.repository.Entity(req.EntityId)
	if err != nil {
		return err
	}

	attr, err := models.NewAttribute(entity, req.Id, req.NameDb, req.TypeDb, req.Default, req.FkTable, req.FkType, req.Nullable, req.PrimaryKey)
	if err != nil {
		return err
	}

	c.repository.CreateAttribute(attr)

	service := entity.Service()

	// todo sql

	return service.Flush()
}

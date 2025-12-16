package entity

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/http/request"
	"github.com/vologzhan/maker-gui/backend/http/response"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
)

type List struct {
	repository *repository.Repository
}

func (c *List) Handle(ctx echo.Context) error {
	var req request.EntityList
	if err := ctx.Bind(&req); err != nil {
		return err
	}

	res, err := c.handle(req)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, res)
}

func (c *List) handle(req request.EntityList) (*response.EntityList, error) {
	service, err := c.repository.Service(req.ServiceId)
	if err != nil {
		return nil, err
	}

	entities, err := service.Entities()
	if err != nil {
		return nil, err
	}

	for _, entity := range entities {
		c.repository.CreateEntity(entity)

		attrs, err := entity.Attributes()
		if err != nil {
			return nil, err
		}

		for _, a := range attrs {
			c.repository.CreateAttribute(a)
		}
	}

	return response.NewEntityList(entities), nil
}

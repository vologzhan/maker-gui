package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker"
	"github.com/vologzhan/maker-gui/backend/http/controller/attribute"
	"github.com/vologzhan/maker-gui/backend/http/controller/entity"
	"github.com/vologzhan/maker-gui/backend/http/controller/service"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
)

func New(server *echo.Echo, mak *maker.Node) {
	rep := repository.New(mak)

	service.New(server, rep)
	entity.New(server, rep)
	attribute.New(server, rep)
}

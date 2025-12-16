package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
)

func New(server *echo.Echo, rep *repository.Repository) {
	g := server.Group("/api/service")

	g.Add(http.MethodGet, "", (&List{rep}).Handle)
	g.Add(http.MethodPost, "", (&Create{rep}).Handle)
}

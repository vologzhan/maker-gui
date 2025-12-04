package entity

import (
	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/repository"
	"net/http"
)

func New(server *echo.Echo, rep *repository.Repository) {
	g := server.Group("/api/entity")

	g.Add(http.MethodGet, "/service/:service-id", (&List{rep}).Handle)
	g.Add(http.MethodPost, "/service/:service-id", (&Create{rep}).Handle)
	g.Add(http.MethodDelete, "/:id", (&Delete{rep}).Handle)
	g.Add(http.MethodPut, "/:id", (&Update{rep}).Handle)
}

package attribute

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vologzhan/maker-gui/backend/maker/repository"
)

func New(server *echo.Echo, rep *repository.Repository) {
	g := server.Group("/api/attribute")

	g.Add(http.MethodPost, "/entity/:entity-id", (&Create{rep}).Handle)
	g.Add(http.MethodDelete, "/:id", (&Delete{rep}).Handle)
	g.Add(http.MethodPut, "/:id", (&Update{rep}).Handle)
}

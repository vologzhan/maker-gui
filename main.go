package main

import (
	"embed"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vologzhan/maker"
	"github.com/vologzhan/maker-gui/backend/controller"
	"github.com/vologzhan/maker/template"
	"log"
	"net/http"
)

//go:embed frontend/dist/*
var frontendFs embed.FS

//go:embed templates/*
var templateFs embed.FS

func main() {
	server := newServer()

	err := newBackend(server)
	if err != nil {
		log.Fatal(err)
	}

	newFrontend(server)

	port := 1551

	err = server.Start(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}

func newBackend(server *echo.Echo) error {
	srcDir, err := getSourceDir()
	if err != nil {
		return err
	}

	mak, err := newMaker(srcDir)
	if err != nil {
		return err
	}

	controller.New(server, mak)

	return nil
}

func newFrontend(server *echo.Echo) {
	subFS := echo.MustSubFS(frontendFs, "frontend/dist")
	fileServer := echo.WrapHandler(http.FileServer(http.FS(subFS)))

	server.GET("/*", func(ctx echo.Context) error {
		return fileServer(ctx)
	})
}

func newMaker(srcDir string) (*maker.Node, error) {
	tpl, err := template.New(templateFs, "templates/current")
	if err != nil {
		return nil, err
	}

	node, err := maker.New(tpl, srcDir)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func newServer() *echo.Echo {
	server := echo.New()
	server.Use(middleware.Logger())
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	return server
}

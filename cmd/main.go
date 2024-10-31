package main

import (
	"ProyectoInge/templates"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

// e.GET("/show", show)
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func main() {
	port := ":3000"
	fmt.Printf("http://localhost%s", port)

	e := echo.New()

	e.Static("/assets", "assets")

	component := templates.Index()

	e.GET("/", func(c echo.Context) error {
		return component.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/uwu", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, uwu!")
	})

	e.GET("/show", show)

	e.Logger.Fatal(e.Start(port))
}

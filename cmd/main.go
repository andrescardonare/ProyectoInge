package main

import (
	"ProyectoInge/templates"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// e.GET("/show", show)
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.QueryParam("team")
	member := c.QueryParam("member")
	return c.String(http.StatusOK, "team:"+team+", member:"+member)
}

func main() {
	loadEnv()
	//controllers.DBconnection()

	port := ":3000"
	fmt.Printf("http://localhost%s", port)

	e := echo.New()

	e.Static("/assets", "assets")
	e.File("/vendored/htmx_v2.0.3.min.js", "vendored/htmx_v2.0.3.min.js")

	home := templates.Index()
	login := templates.Login()

	e.GET("/", func(c echo.Context) error {
		return home.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/tracking", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, uwu!")
	})

	e.GET("/login", func(c echo.Context) error {
		return login.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "logout")
	})

	e.GET("/app", func(c echo.Context) error {
		return c.String(http.StatusOK, "app")
	})

	e.GET("/show", show)

	e.Logger.Fatal(e.Start(port))
}

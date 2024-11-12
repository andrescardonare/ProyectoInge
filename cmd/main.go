package main

import (
	"ProyectoInge/controllers"
	"ProyectoInge/templates"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
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
	//controllers.register()

	if len(os.Args) < 1 {
		controllers.DBconnection()
	} else {
		fmt.Println("Modo Demo")
	}

	port := ":3000"
	fmt.Printf("http://localhost%s", port)

	e := echo.New()

	e.Static("/assets", "assets")
	e.File("/vendored/htmx_v2.0.3.min.js", "vendored/htmx_v2.0.3.min.js")

	e.Use()

	home := templates.Index()
	auth := templates.Auth()
	register := templates.Register()
	login := templates.LogIn()
	trck := templates.Tracking()

	e.GET("/", func(c echo.Context) error {
		return home.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/tracking", func(c echo.Context) error {
		return trck.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/auth", func(c echo.Context) error {
		return auth.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/register", func(c echo.Context) error {
		return register.Render(context.Background(), c.Response().Writer)
	})

	e.POST("/register", controllers.Register)

	e.GET("/login", func(c echo.Context) error {
		return login.Render(context.Background(), c.Response().Writer)
	})

	e.POST("/login", controllers.Login)

	e.GET("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "logout")
	})

	e.GET("/app", func(c echo.Context) error {
		return c.String(http.StatusOK, "app")
	})

	e.GET("/protected", controllers.Protected, controllers.AuthMiddleware)

	e.GET("/show", show)

	e.Logger.Fatal(e.Start(port))
}

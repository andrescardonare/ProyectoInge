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
	if len(os.Args) < 1 {
		loadEnv()
		controllers.DBconnection()
	} else {
		fmt.Println("Modo Demo")
	}

	port := ":3000"
	fmt.Printf("http://localhost%s", port)

	e := echo.New()

	e.Static("/assets", "assets")
	e.File("/vendored/htmx_v2.0.3.min.js", "vendored/htmx_v2.0.3.min.js")
	e.File("/DEMO_MODE", "templates/tables.html")

	/*
		e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
			// Be careful to use constant time comparison to prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(username), []byte("demo_utb")) == 1 &&
				subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
				return true, nil
			}
			return false, nil
		}))
	*/

	home := templates.Index()
	auth := templates.Auth()
	register := templates.Register()
	login := templates.LogIn()
	trck := templates.Tracking()
	acknow := templates.Acknowledgement()
	//protec := templates.Protected()

	e.GET("/", func(c echo.Context) error {
		return home.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/tracking", func(c echo.Context) error {
		return trck.Render(context.Background(), c.Response().Writer)
	})

	e.GET("/acknowledgements", func(c echo.Context) error {
		return acknow.Render(context.Background(), c.Response().Writer)
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

	e.POST("/login", func(c echo.Context) error {
		return c.String(http.StatusOK, "Deshabilitado para Demostración")
	})

	e.GET("/logout", func(c echo.Context) error {
		return c.String(http.StatusOK, "logout")
	})

	e.GET("/app", func(c echo.Context) error {
		return c.File("tables.html")
	})

	//e.GET("/protected")

	e.GET("/show", show)

	e.POST("/submit-form", controllers.PostToApi)

	e.GET("/map", func(c echo.Context) error {
		newImageHTML := `
			<img 
				src="/assets/images/mapita.png" 
				alt="New Placeholder Image" 
				class="w-[48rem] max-w-none rounded-xl shadow-xl ring-1 ring-gray-400/10 sm:w-[57rem]">
		`
		return c.HTML(http.StatusOK, newImageHTML)
	})

	e.Logger.Fatal(e.Start(port))
}

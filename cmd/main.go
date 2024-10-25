package main

import (
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

	e.Static("/static", "assets")

	//component := templates.Hello("John")
	//component.Render(context.Background(), os.Stdout)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/uwu", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, uwu!")
	})

	e.GET("/show", show)

	e.Logger.Fatal(e.Start(port))
}

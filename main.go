package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		tmpl, _ := template.ParseFiles("templates/index.html")
		tmpl.ExecuteTemplate(w, "index.html", nil)
	})
	http.ListenAndServe(":3000", r)
}

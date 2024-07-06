package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// GET https://api.spoonacular.com/recipes/complexSearch
func getRecipes(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", getRecipes)
	http.ListenAndServe(":3000", r)
}

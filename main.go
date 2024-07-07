package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

type RecipeResponse struct {
	Offset       int `json:"offset"`
	Number       int `json:"number"`
	TotalResults int `json:"totalResults"`
	Results      []struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
		Image string `json:"image"`
	} `json:"results"`
}

// GET https://api.spoonacular.com/recipes/complexSearch?apiKey=97ef6c2a62534b28b84d57045ada4a52&query=pasta&number=4

func getRecipes(query string, number int) (*RecipeResponse, error) {
	apiKey := os.Getenv("SPOONACULAR_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("SPOONACULAR_API_KEY environment variable is not set")
	}

	url := fmt.Sprintf("https://api.spoonacular.com/recipes/complexSearch?apiKey=%s&query=%s&number=%d", apiKey, query, number)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var recipesResponse RecipeResponse
	if err := json.NewDecoder(resp.Body).Decode(&recipesResponse); err != nil {
		return nil, err
	}

	return &recipesResponse, nil
}

func recipesReporter(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	numberStr := r.URL.Query().Get("number")

	if query == "" || numberStr == "" {
		http.Error(w, "Query and number parameters are required", http.StatusBadRequest)
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		http.Error(w, "Invalid number parameter", http.StatusBadRequest)
		return
	}

	recipes, err := getRecipes(query, number)
	if err != nil {
		log.Printf("Could not get recipes: %v", err)
		http.Error(w, "Could not get recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	renderTemplate(w, recipes, "food-app.html")
}

func renderTemplate(w http.ResponseWriter, recipe *RecipeResponse, file string) {
	tmpl, err := template.ParseFiles(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, recipe); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", recipesReporter)
	http.ListenAndServe(":3000", r)
}

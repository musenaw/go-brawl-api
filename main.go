package main

import (
	"fmt"
	"musenaw/go-brawl-api/controllers"
	"musenaw/go-brawl-api/models"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	cfg := models.DefaultPostgresConfig()
	err := models.Open(cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	r.Get("/", controllers.StaticHandlerJSON)
	r.Get("/players/{playerId}", controllers.GetPlayerInfo)
	r.Get("/players/{playerId}/battlelog", controllers.GetPlayerBattlelog)

	r.Post("/signup", controllers.StaticHandlerJSON)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}

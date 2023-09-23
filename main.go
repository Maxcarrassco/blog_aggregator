package main

import (
	"fmt"
	"net/http"

	"github.com/Maxcarrassco/blog_aggregator/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)


func main() {
	app := chi.NewRouter()
	app.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge: 300,
	}))
	app.Get("/api/v1/readiness", func(w http.ResponseWriter, r *http.Request) {
		msg := struct {
			Status string `json:"status"`
		}{ Status: "ok" }
		utils.ResponseWithJSON(w, 200, msg)
	})
	app.Get("/api/v1/err", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseWithError(w, 500, "internal server error")
	})
	const PORT = "8080"
	addr := fmt.Sprintf("127.0.0.1:%s", PORT)
	fmt.Printf("Server is listening on port %s!", PORT)
	http.ListenAndServe(addr, app)
}

package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/Maxcarrassco/blog_aggregator/internal/database"
	"github.com/Maxcarrassco/blog_aggregator/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func (self *APIConfig) createUser(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseWithError(w, 400, "invalid request body")
		return
	}
	body := struct {
		Name string `json:"name"`
	}{}
	err = json.Unmarshal(data, &body)
	if err != nil || body.Name == ""{
		utils.ResponseWithError(w, 400, "invalid request body")
		return
	}
	today := time.Now()
	id := uuid.NewString()
	api_key := fmt.Sprintf("%x", sha256.Sum256([]byte(id)))
	userParam := database.CreateUserParams{
		ID: id,
		Name: body.Name,
		ApiKey: api_key,
		CreatedAt: today,
		UpdatedAt: sql.NullTime{Time: today, Valid: true},
	}
	_, err = self.DB.CreateUser(self.Ctx, userParam)
	if err != nil {
		utils.ResponseWithError(w, 400, "Unable to create user")
		return
	}
	res := struct {
		ID string `json:"id"`
		Name string `json:"name"`
		ApiKey string `json:"api_key"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} {ID: userParam.ID, Name: userParam.Name, ApiKey: api_key, CreatedAt: userParam.CreatedAt, UpdatedAt: userParam.UpdatedAt.Time,}
	utils.ResponseWithJSON(w, 201, res)
}


func (self *APIConfig) getUserByApiKey(w http.ResponseWriter, r *http.Request) {
	authz := r.Header.Get("Authorization")
	if authz == "" {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	auth := strings.Split(authz, " ")
	if len(auth) != 2 || auth[0] != "ApiKey" || auth[1] == "" {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}

	userParam, err := self.DB.GetUserByApiKey(self.Ctx, auth[1])
	if err != nil {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	res := struct {
		ID string `json:"id"`
		Name string `json:"name"`
		ApiKey string `json:"api_key"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} {ID: userParam.ID, Name: userParam.Name, ApiKey: userParam.ApiKey, CreatedAt: userParam.CreatedAt, UpdatedAt: userParam.UpdatedAt.Time,}
	utils.ResponseWithJSON(w, 200, res)
}


type APIConfig struct {
	DB *database.Queries
	Ctx context.Context
}


func New() (*APIConfig, error) {
	defer func() {
        if r := recover(); r != nil {
            log.Println("Stack Track\n", string(debug.Stack()))
        }
    	}()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return &APIConfig{}, err
	}
	ctx := context.Background()
	DB_URL := os.Getenv("DB_URL")
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatal("Error etsablishing db cnnection: ", err.Error())
		return &APIConfig{}, err
	}
	dbQueries := database.New(db)
	return &APIConfig{DB: dbQueries, Ctx: ctx}, nil
}


func main() {
	app := chi.NewRouter()
	app.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"POST", "GET", "PUT", "OPTIONS", "DELETE"},
		AllowedHeaders: []string{"Authorization", "Content-Type", "X-CSRF-Token"},
		MaxAge: 300,
	}))
	defer func() {
        if r := recover(); r != nil {
            log.Println("Stack Track\n", string(debug.Stack()))
        }
    	}()
	apiConfig, err := New()
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	app.Get("/api/v1/readiness", func(w http.ResponseWriter, r *http.Request) {
		msg := struct {
			Status string `json:"status"`
		}{ Status: "ok" }
		utils.ResponseWithJSON(w, 200, msg)
	})
	app.Get("/api/v1/err", func(w http.ResponseWriter, r *http.Request) {
		utils.ResponseWithError(w, 500, "Internal Server Error")
	})
	app.Post("/api/v1/users", apiConfig.createUser)
	app.Get("/api/v1/users", apiConfig.getUserByApiKey)
	const PORT = "8080"
	addr := fmt.Sprintf("127.0.0.1:%s", PORT)
	fmt.Printf("Server is listening on port %s!", PORT)
	http.ListenAndServe(addr, app)
}

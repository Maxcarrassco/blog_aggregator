package handler

import (
	"context"
	"database/sql"
	"log"
	"os"
	"runtime/debug"

	"github.com/Maxcarrassco/blog_aggregator/internal/database"
	"github.com/joho/godotenv"
)




type APIConfig struct {
	DB *database.Queries
	Ctx context.Context
}


func new() *APIConfig {
	defer func() {
        if r := recover(); r != nil {
            log.Println("Stack Track\n", string(debug.Stack()))
        }
    	}()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		panic(err)
	}
	ctx := context.Background()
	DB_URL := os.Getenv("DB_URL")
	db, err := sql.Open("mysql", DB_URL)
	if err != nil {
		log.Fatal("Error etsablishing db cnnection: ", err.Error())
		panic(err)
	}
	dbQueries := database.New(db)
	return &APIConfig{DB: dbQueries, Ctx: ctx}
}


var ApiCfg = new()

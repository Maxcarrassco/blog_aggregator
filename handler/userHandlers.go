package handler

import (
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"
	"net/http"

	"github.com/Maxcarrassco/blog_aggregator/internal/database"
	"github.com/Maxcarrassco/blog_aggregator/utils"
	"github.com/google/uuid"
)


func (self APIConfig) CreateUser(w http.ResponseWriter, r *http.Request) {
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


func (self APIConfig) GetUserByApiKey(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user")
	if user == nil {
		utils.ResponseWithError(w, 401, "unauthoized")
		return;
	}
	userParam := user.(database.User)
	res := struct {
		ID string `json:"id"`
		Name string `json:"name"`
		ApiKey string `json:"api_key"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} {ID: userParam.ID, Name: userParam.Name, ApiKey: userParam.ApiKey, CreatedAt: userParam.CreatedAt, UpdatedAt: userParam.UpdatedAt.Time,}
	utils.ResponseWithJSON(w, 200, res)
}

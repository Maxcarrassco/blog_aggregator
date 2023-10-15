package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Maxcarrassco/blog_aggregator/internal/database"
	"github.com/Maxcarrassco/blog_aggregator/utils"
	"github.com/google/uuid"
)



func (self APIConfig) CreateFeed(w http.ResponseWriter, r *http.Request) {
	userCxt := r.Context().Value("user")
	if userCxt == nil {
		utils.ResponseWithError(w, 401, "unauthorized")
		return
	}
	user := userCxt.(database.User)

	body, err := io.ReadAll(r.Body)

	if err != nil {
		utils.ResponseWithError(w, 401, "invalid request")
		return
	}

	data := struct {
		Name string `json:"name"`
		Url string `json:"url"`
	}{}

	err = json.Unmarshal(body, &data)

	if err != nil || data.Name == "" || data.Url == "" {
		utils.ResponseWithError(w, 401, "invalid request")
		return
	}

	today := time.Now()
	id := uuid.NewString()

	feedParam := database.CreateFeedParams{
		Name: data.Name,
		ID: id,
		Url: data.Url,
		UserID: user.ID,
		CreatedAt: today,
		UpdatedAt: sql.NullTime{Time: today, Valid: true},
	}

	_, err = self.DB.CreateFeed(self.Ctx, feedParam)
	if err != nil {
		utils.ResponseWithError(w, 400, "unable to create feed")
		return
	}

	res := struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Url string `json:"url"`
		UserID string `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} {ID: feedParam.ID, Name: feedParam.Name, Url: feedParam.Url, UserID: feedParam.UserID, CreatedAt: feedParam.CreatedAt, UpdatedAt: feedParam.UpdatedAt.Time, }
	utils.ResponseWithJSON(w, 201, res)
}

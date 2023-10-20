package handler

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Maxcarrassco/blog_aggregator/internal/database"
	"github.com/Maxcarrassco/blog_aggregator/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)



func (self APIConfig) CreateFeedFollows(w http.ResponseWriter, r *http.Request) {
	userDB := r.Context().Value("user")
	if userDB == nil {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	user, ok := userDB.(database.User)
	if !ok {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.ResponseWithError(w, 400, "invalid request body")
		return
	}
	data := struct {
		FeedId string `json:"feed_id"`
	}{}
	err = json.Unmarshal(body, &data)
	if err != nil || data.FeedId == "" {
		utils.ResponseWithError(w, 400, "invalid request body")
		return
	}
	_, err = self.DB.GetFeedById(self.Ctx, data.FeedId)
	if err != nil {
		utils.ResponseWithError(w, 404, "feed does not exists")
		return
	}
	today := time.Now()
	id := uuid.NewString()
	FeedFollows := database.CreateFeedFollowParams{
		ID: id,
		FeedID: data.FeedId,
		UserID: user.ID,
		CreatedAt: today,
		UpdatedAt: sql.NullTime{Time: today, Valid: true, },
	}
	_, err = self.DB.CreateFeedFollow(self.Ctx, FeedFollows)
	if err != nil {
		utils.ResponseWithError(w, 400, "unable to create feed follow")
		return
	}
	res := struct {
		Id string `json:"id"`
		FeedId string `json:"feed_id"`
		UserId string `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} {Id: FeedFollows.ID, FeedId: FeedFollows.FeedID, UserId: FeedFollows.UserID, CreatedAt: FeedFollows.CreatedAt, UpdatedAt: FeedFollows.UpdatedAt.Time,}
	utils.ResponseWithJSON(w, 201, res)
}

func (self APIConfig) GetUserFeedFollows(w http.ResponseWriter, r *http.Request) {
	userDB := r.Context().Value("user")
	if userDB == nil {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	user, ok := userDB.(database.User)
	if !ok {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	result, err := self.DB.GetUserFeedFollow(self.Ctx, user.ID)
	if err != nil {
		utils.ResponseWithError(w, 400, "unable to get user feed follows")
		return
	}
	type res struct {
		Id string `json:"id"`
		FeedId string `json:"feed_id"`
		UserId string `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	feeds := []res{}

	for _, FeedFollows := range result {
		feed := res {Id: FeedFollows.ID, FeedId: FeedFollows.FeedID, UserId: FeedFollows.UserID, CreatedAt: FeedFollows.CreatedAt, UpdatedAt: FeedFollows.UpdatedAt.Time,}
		feeds = append(feeds, feed)
	}
	utils.ResponseWithJSON(w, 200, feeds)
}

func (self APIConfig) DeleteFeedFollow(w http.ResponseWriter, r *http.Request) {
	userDB := r.Context().Value("user")
	if userDB == nil {
		utils.ResponseWithError(w, 401, "unauthorized request")
		return
	}
	id := chi.URLParam(r, "id")
	FeedFollow, err := self.DB.GetFeedFollowById(self.Ctx, id)
	if err != nil {
		utils.ResponseWithError(w, 404, "feed follow not found")
		return
	}
	_, err = self.DB.DeleteUserFeedFollow(self.Ctx, id)
	if err != nil {
		utils.ResponseWithError(w, 400, "unable to delete user feed follows")
		return
	}
	res := struct {
		Id string `json:"id"`
		FeedId string `json:"feed_id"`
		UserId string `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} {Id: FeedFollow.ID, FeedId: FeedFollow.FeedID, UserId: FeedFollow.UserID, CreatedAt: FeedFollow.CreatedAt, UpdatedAt: FeedFollow.UpdatedAt.Time,}
	utils.ResponseWithJSON(w, 200, res)
}

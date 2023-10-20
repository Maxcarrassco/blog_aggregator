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
	fId := uuid.NewString()
	feedFollowParams := database.CreateFeedFollowParams{
		ID: fId,
		UserID: user.ID,
		FeedID: feedParam.ID,
		CreatedAt: today,
		UpdatedAt: sql.NullTime{Time: today, Valid: true, },
	}

	self.DB.CreateFeedFollow(self.Ctx, feedFollowParams)

	type feedS struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Url string `json:"url"`
		UserID string `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} 
	feedRes := feedS {ID: feedParam.ID, Name: feedParam.Name, Url: feedParam.Url, UserID: feedParam.UserID, CreatedAt: feedParam.CreatedAt, UpdatedAt: feedParam.UpdatedAt.Time, }
	type feedFellowS struct {
		ID string `json:"id"`
		UserID string `json:"user_id"`
		FeedID string `json:"feed_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} 
	feedFellowRes := feedFellowS {ID: feedFollowParams.ID, UserID: feedFollowParams.UserID, FeedID: feedFollowParams.FeedID, CreatedAt: feedFollowParams.CreatedAt, UpdatedAt: feedFollowParams.UpdatedAt.Time, }
	res := struct {
		Feed feedS `json:"feed"`
		FeedFollow feedFellowS `json:"feed_follow"`
	} {Feed: feedRes, FeedFollow: feedFellowRes, }
	utils.ResponseWithJSON(w, 201, res)
}


func (self APIConfig) GetAllFeeds(w http.ResponseWriter, r *http.Request) {
	type res struct {
		ID string `json:"id"`
		Name string `json:"name"`
		Url string `json:"url"`
		UserID string `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
	out := []res{}
	feeds, err := self.DB.GetAllFeeds(self.Ctx)
	if err != nil {
		utils.ResponseWithError(w, 400, "unable to get feeds")
		return
	}

	for _, feedParam := range feeds {
		data := res {ID: feedParam.ID, Name: feedParam.Name, Url: feedParam.Url, UserID: feedParam.UserID, CreatedAt: feedParam.CreatedAt, UpdatedAt: feedParam.UpdatedAt.Time, }
		out = append(out, data)
	}
	utils.ResponseWithJSON(w, 200, out)
}

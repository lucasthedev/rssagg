package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/lucasthedev/rssagg/internal/database"
)

func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Erro ao fazer o parse do json: %v", err))
		return
	}

	feedFollow, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Erro ao criar feed follow: %v", err))
		return
	}

	respondWithJson(w, 201, dataBaseFeedFollowToFeedFollow(feedFollow))
}

func (apiCfg *apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Erro buscar os feed follows: %v", err))
		return
	}

	respondWithJson(w, 201, dataBaseFeedFollowsToFeedsFollows(feedFollows))
}

func (apiCfg *apiConfig) handlerUnfollowFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowId := chi.URLParam(r, "feedFollowID")
	feedFllowId, err := uuid.Parse(feedFollowId)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Erro ao fazer o parse do uuid: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFllowId,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Erro ao dar unfollow feed: %v", err))
		return
	}

	respondWithJson(w, 200, struct{}{})
}

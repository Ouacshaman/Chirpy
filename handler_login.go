package main

import (
	"Chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "Unable to Get User", err)
		return
	}
	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Unauthorized", err)
		return
	}
	var token string
	if params.ExpiresInSeconds == 0 ||
		params.ExpiresInSeconds > 3600 {
		hour, err := time.ParseDuration("1h")
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't Parse Time", err)
			return
		}
		outcome, err := auth.MakeJWT(user.ID, cfg.secret, hour)
		if err != nil {
			respondWithError(w, 401, "Unauthorized", err)
			return
		}
		token = outcome
	} else {
		stringSeconds := strconv.Itoa(params.ExpiresInSeconds) + "s"
		secs, err := time.ParseDuration(stringSeconds)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Can't Parse Time", err)
			return
		}
		outcome, err := auth.MakeJWT(user.ID, cfg.secret, secs)
		if err != nil {
			respondWithError(w, 401, "Unauthorized", err)
			return
		}
		token = outcome
	}
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	respondWithJSON(w, 200, returnVals{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})
}

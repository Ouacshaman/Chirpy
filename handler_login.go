package main

import (
	"Chirpy/internal/auth"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	type returnVals struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	respondWithJSON(w, 200, returnVals{
		Id:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

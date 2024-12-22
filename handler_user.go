package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 404, "Unable to Hash Password", err)
		return
	}
	userParam := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPw,
	}
	user, err := cfg.db.CreateUser(r.Context(), userParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to Create User", err)
		return
	}
	fmt.Println("New User UUID:", user.ID)
	type returnVals struct {
		Id          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	respondWithJSON(w, 201, returnVals{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

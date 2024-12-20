package main

import (
	"Chirpy/internal/auth"
	"Chirpy/internal/database"
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w,
			401,
			"Unable to obtain bearer token",
			err)
		return
	}
	userID, err := auth.ValidateJWT(bt, cfg.secret)
	if err != nil {
		respondWithError(w,
			401,
			"Unable to validate token",
			err)
		return
	}
	type paramaters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w,
			401,
			"Couldn't decode parameters", err)
		return
	}
	hashed, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w,
			401,
			"Couldn't hash password", err)
		return
	}
	err = cfg.db.UpdateUserEmailPw(r.Context(),
		database.UpdateUserEmailPwParams{
			HashedPassword: hashed,
			Email:          params.Email,
			ID:             userID,
		})
	if err != nil {
		respondWithError(w,
			401,
			"Couldn't make an update", err)
		return
	}
	type returnVals struct {
		Email string `json:"email"`
	}
	respondWithJSON(w, 200,
		returnVals{
			Email: params.Email,
		})
}

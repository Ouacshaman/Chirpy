package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Id        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}
	s := r.URL.Query().Get("author_id")
	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 401, "unable to parse string", err)
			return
		}
		chirps, err := cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, 400, "Unable To Get Chirps", err)
			return
		}
		returnVal := []chirp{}
		for _, v := range chirps {
			tempChirp := chirp{
				Id:        v.ID,
				CreatedAt: v.CreatedAt,
				UpdatedAt: v.UpdatedAt,
				Body:      v.Body,
				UserId:    v.UserID,
			}
			returnVal = append(returnVal, tempChirp)
		}
		respondWithJSON(w, 200, returnVal)
		return
	}
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Unable To Get Chirps", err)
		return
	}
	returnVal := []chirp{}
	for _, v := range chirps {
		tempChirp := chirp{
			Id:        v.ID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Body:      v.Body,
			UserId:    v.UserID,
		}
		returnVal = append(returnVal, tempChirp)
	}
	respondWithJSON(w, 200, returnVal)
}

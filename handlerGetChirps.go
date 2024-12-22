package main

import (
	"Chirpy/internal/database"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	if s != "" {
		authorID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 401, "unable to parse string", err)
			return
		}
		order := r.URL.Query().Get("sort")
		if order == "desc" {
			chirps, err := cfg.db.GetChirpsByUserIDDesc(r.Context(), authorID)
			if err != nil {
				respondWithError(w, 400, "Unable To Get Chirps", err)
				return
			}
			returnVal(w, chirps)
			return
		}
		chirps, err := cfg.db.GetChirpsByUserID(r.Context(), authorID)
		if err != nil {
			respondWithError(w, 400, "Unable To Get Chirps", err)
			return
		}
		returnVal(w, chirps)
		return
	}
	order := r.URL.Query().Get("sort")
	if order == "desc" {
		chirps, err := cfg.db.GetChirpsDesc(r.Context())
		if err != nil {
			respondWithError(w, 400, "Unable To Get Chirps", err)
			return
		}
		returnVal(w, chirps)
		return

	}
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 400, "Unable To Get Chirps", err)
		return
	}
	returnVal(w, chirps)
}

func returnVal(w http.ResponseWriter, chirps []database.Chirp) {
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

package main

import (
	"Chirpy/internal/auth"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	bt, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w,
			401,
			"Unable to obtain bearer token",
			err)
		return
	}
	vid, err := auth.ValidateJWT(bt, cfg.secret)
	if err != nil {
		respondWithError(w,
			403,
			"Unable to validate token",
			err)
		return
	}
	chirpUUID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, 404, "Unable to Convert String to UUID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, 401,
			"Unable to Get Chirp", err)
		return
	}
	if chirp.UserID != vid {
		err := errors.New("User Not Authorized")
		respondWithError(w, 403,
			"Non-Authorized User", err)
		return
	}
	err = cfg.db.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, 404,
			"Unable to Delete Chirp", err)
		return
	}
	respondWithJSON(w, 204, struct{}{})
}

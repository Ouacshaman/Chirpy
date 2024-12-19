package main

import (
	"Chirpy/internal/auth"
	"net/http"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	rtk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w,
			http.StatusBadRequest,
			"Unable to Generate Refresh Token", err)
	}
	err = cfg.db.RevokeToken(r.Context(), rtk)
	if err != nil {
		respondWithError(w,
			http.StatusBadRequest,
			"Unable to Revoke Token", err)
	}
	respondWithJSON(w, 204, struct{}{})
}

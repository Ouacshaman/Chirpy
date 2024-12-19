package main

import (
	"Chirpy/internal/auth"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	rtk, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w,
			http.StatusBadRequest,
			"Unable to Generate Refresh Token", err)
	}
	grt, err := cfg.db.GetRefreshToken(r.Context(),
		rtk)
	if err != nil || grt.ExpiresAt == time.Now() {
		respondWithError(w,
			401,
			"Unable to Retrieve UserID", err)
	}
	nt, err := auth.MakeJWT(grt.UserID, rtk)
	if err != nil {
		respondWithError(w,
			http.StatusBadRequest,
			"Unable to Generate Access Token", err)
	}
	type returnVals struct {
		Token string `json:"token"`
	}
	respondWithJSON(w, 200, returnVals{
		Token: nt,
	})
}

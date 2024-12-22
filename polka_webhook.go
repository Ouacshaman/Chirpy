package main

import (
	"Chirpy/internal/auth"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "ApiKey Not Found", err)
		return
	}
	if cfg.polkaKey != apiKey {
		respondWithError(w, 401, "Invalid ApiKey", err)
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, struct{}{})
		return
	}
	err = cfg.db.UpgradeChripy(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, 404, "Not Found", err)
	}
	respondWithJSON(w, 204, struct{}{})
}

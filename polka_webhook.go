package main

import (
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
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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

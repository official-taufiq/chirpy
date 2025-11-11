package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/official-taufiq/chirpy/internal/auth"
)

func (cfg apiConfig) chirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		}
	}

	apikey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Malformed or no api key provided")
		return
	}
	if apikey != cfg.apikey {
		respondWithError(w, 401, "Invalid api key")
		return
	}

	params := parameters{}

	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode body")
		return
	}
	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	_, err = cfg.db.ChirpRed(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(404)
		return
	}
	w.WriteHeader(204)

}

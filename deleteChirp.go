package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/official-taufiq/chirpy/internal/auth"
)

func (cfg apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")

	chirpId, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp Id")
		return
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Malformed or no bearer token provided")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid bearer token")
		return
	}

	dbChirp, err := cfg.db.GetOneChirp(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, "no chirp found")
		return
	}
	if dbChirp.UserID == userId {
		err = cfg.db.DeleteChirp(r.Context(), chirpId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "couldn't delete chirp")
			return
		}
		w.WriteHeader(204)
		return
	}
	respondWithError(w, 403, "Not authorized to delete another user's chirp")
}

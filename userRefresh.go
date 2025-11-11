package main

import (
	"net/http"
	"time"

	"github.com/official-taufiq/chirpy/internal/auth"
)

func (cfg apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No bearer token provided")
		return
	}
	dbUser, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.tokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't create access token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg apiConfig) tokenRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No bearer token provided")
		return
	}
	_, err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke session")
		return
	}
	w.WriteHeader(204)

}

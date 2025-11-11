package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/official-taufiq/chirpy/internal/auth"
	"github.com/official-taufiq/chirpy/internal/database"
)

type chirps struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}
type errorResponse struct {
	Error string `json:"error"`
}

func (cfg apiConfig) chirp(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find bearer token")
		return
	}
	userId, err := auth.ValidateJWT(token, cfg.tokenSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't validate bearer token")
		return
	}
	decoder := json.NewDecoder(r.Body)
	type params struct {
		Body string `json:"Body"`
	}
	req := params{}
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, 405, "Something went wrong")
		return
	}
	if len(req.Body) > 140 {
		respondWithError(w, 400, "Chirp too long")
		return
	}
	clean := profanityChecker(req.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   clean,
		UserID: userId,
	})
	if err != nil {
		log.Printf("couldn't create chirp: %v", err)
		respondWithError(w, http.StatusInternalServerError, "coundn't create chirp")
		return
	}
	respondWithJSON(w, http.StatusOK, chirps{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func profanityChecker(text string) string {
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(text, " ")

	for i, word := range words {
		for _, prof := range profaneWords {
			if strings.ToLower(word) == prof {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}

func (cfg apiConfig) AllChirps(w http.ResponseWriter, r *http.Request) {

	res, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("error getting allChirps: %v", err)

		respondWithError(w, http.StatusInternalServerError, "couldn't fetch all chirps")
		return
	}
	authorID := uuid.Nil
	IDstring := r.URL.Query().Get("author_id")

	if IDstring != "" {
		authorID, err = uuid.Parse(IDstring)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID")
			return
		}
	}

	sortDirection := "asc"
	sortDirectionParams := r.URL.Query().Get("sort")
	if sortDirectionParams == "desc" {
		sortDirection = "desc"
	}

	Chirps := []chirps{}
	for _, dbChirp := range res {
		if authorID != uuid.Nil && dbChirp.UserID != authorID {
			continue
		}

		Chirps = append(Chirps, chirps{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		})
	}

	sort.Slice(Chirps, func(i, j int) bool {
		if sortDirection == "desc" {
			return Chirps[i].CreatedAt.After(Chirps[j].CreatedAt)
		}
		return Chirps[i].CreatedAt.Before(Chirps[j].CreatedAt)
	})
	respondWithJSON(w, 200, Chirps)
}

func (cfg apiConfig) OneChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDstring := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIDstring)
	if err != nil {
		log.Printf("invalid chirp ID: %v", err)
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}
	dbChirp, err := cfg.db.GetOneChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("chirp with this ID doesn't exist:%v", err)
		respondWithError(w, http.StatusBadRequest, "Chirp doesn't exist")
		return
	}
	respondWithJSON(w, 200, chirps{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	})
}

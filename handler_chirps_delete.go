package main

import (
	"net/http"

	"github.com/brettlazarine/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	if chirpID == "" {
		respondWithError(w, http.StatusNotFound, "missing id parameter", nil)
		return
	}
	uuidChirpID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID format", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "could not get token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), uuidChirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "could not get chirp", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "user does not own chirp", nil)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), uuidChirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not delete chirp", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

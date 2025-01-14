package main

import (
	"net/http"
	"sort"

	"github.com/brettlazarine/Chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	var dbChirps []database.Chirp
	var err error
	if author_id != "" {
		authorId, err := uuid.Parse(author_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing author id", err)
			return
		}
		dbChirps, err = cfg.db.GetChirpsByAuthor(r.Context(), authorId)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps", err)
			return
		}
	} else {
		dbChirps, err = cfg.db.GetAllChirps(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "could not get chirps", err)
			return
		}
	}

	chirps := make([]Chirp, len(dbChirps))

	for i, chirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
	}

	sortType := r.URL.Query().Get("sort")
	if sortType == "asc" || sortType == "" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		})
	} else if sortType == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing chirp id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found", err)
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

// *** Boot.dev Implementation
// func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
// 	dbChirps, err := cfg.db.GetChirps(r.Context())
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps", err)
// 		return
// 	}

// 	authorID := uuid.Nil
// 	authorIDString := r.URL.Query().Get("author_id")
// 	if authorIDString != "" {
// 		authorID, err = uuid.Parse(authorIDString)
// 		if err != nil {
// 			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
// 			return
// 		}
// 	}

// 	chirps := []Chirp{}
// 	for _, dbChirp := range dbChirps {
// 		if authorID != uuid.Nil && dbChirp.UserID != authorID {
// 			continue
// 		}

// 		chirps = append(chirps, Chirp{
// 			ID:        dbChirp.ID,
// 			CreatedAt: dbChirp.CreatedAt,
// 			UpdatedAt: dbChirp.UpdatedAt,
// 			UserID:    dbChirp.UserID,
// 			Body:      dbChirp.Body,
// 		})
// 	}

// 	respondWithJSON(w, http.StatusOK, chirps)
// }

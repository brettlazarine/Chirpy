package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not decode parameters", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: replaceBadWord(params.Body)})
}

func replaceBadWord(body string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	chirpWords := strings.Split(body, " ")

	for i, word := range chirpWords {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				chirpWords[i] = "****"
			}
		}
	}

	return strings.Join(chirpWords, " ")
}

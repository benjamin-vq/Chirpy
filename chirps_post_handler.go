package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

type errorResponse struct {
	Error string `json:"error,omitempty"`
}

func (cfg *apiConfig) postChirpHandler(w http.ResponseWriter, r *http.Request) {

	type chirpParams struct {
		Body string `json:"body"`
	}
	params := chirpParams{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil {
		log.Printf("Error decoding chirp: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode chirp")
		return
	}

	sanitized, err := validateChirp(params.Body)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	chirp, err := cfg.DB.CreateChirp(sanitized)
	if err != nil {
		log.Printf("Could not save chirp to database: %q", err)
		respondWithError(w, 500, "Could not post chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func validateChirp(body string) (string, error) {
	const limit = 140

	if chirpLen := len(body); chirpLen > limit {
		log.Printf("Decoded chirp body (%d) is greater than the limit (%d)", chirpLen, limit)
		return "", errors.New("chirp length exceeds limit")
	}

	sanitized := replaceBadWords(body)
	log.Print("Chirp validated successfully")

	return sanitized, nil
}

func replaceBadWords(body string) (replaced string) {
	oldnew := []string{
		"Kerfuffle", "****",
		"kerfuffle", "****",
		"Sharbert", "****",
		"sharbert", "****",
		"Fornax", "****",
		"fornax", "****",
	}

	r := strings.NewReplacer(oldnew...)

	return r.Replace(body)
}

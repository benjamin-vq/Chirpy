package main

import (
	"encoding/json"
	"errors"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postChirpHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	token, found := strings.CutPrefix(authHeader, "Bearer ")
	if !found {
		log.Printf("Invalid authorization header: %q", authHeader)
		respondWithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	userId, err := auth.UserIdFromToken(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error retrieving user id from token: %q", err)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	type chirpParams struct {
		Body string `json:"body"`
	}
	params := chirpParams{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)

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

	chirp, err := cfg.DB.CreateChirp(sanitized, userId)
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

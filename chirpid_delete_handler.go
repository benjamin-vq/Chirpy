package main

import (
	"errors"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (cfg *apiConfig) deleteChirpIdHandler(w http.ResponseWriter, r *http.Request) {

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

	pv := r.PathValue("chirpId")
	chirpId, err := strconv.Atoi(pv)
	if err != nil {
		log.Printf("Provided chirp id to delete is not valid: %q", err)
		respondWithError(w, http.StatusBadRequest, "Invalid chirp id")
		return
	}

	err = cfg.DB.DeleteChirpById(chirpId, userId)
	if err != nil {
		if errors.Is(err, database.IncorrectAuthorId) || errors.Is(err, database.ChirpNotExists) {
			log.Printf("Received chirp id is incorrect: %q", err)
			respondWithError(w, http.StatusForbidden, "You are not authorized to do that")
			return
		}
		log.Printf("Error received trying to delete chirp: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Internal error")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

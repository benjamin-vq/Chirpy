package main

import (
	"github.com/benjamin-vq/chirpy/internal/auth"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postRefreshHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	refreshToken, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found || refreshToken == "" {
		log.Printf("Did not find a valid authorization header")
		respondWithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	userId, err := cfg.DB.UserIdFromRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Token does not exist or is expired")
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
		return
	}

	newToken, err := auth.CreateJwt(userId, cfg.jwtSecret)
	if err != nil {
		log.Printf("Could not create a new token based on a refresh token: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not refresh token")
		return
	}

	type response struct {
		NewToken string `json:"token"`
	}

	log.Printf("Successfully generated a new token based on a refresh token for user with id: %d", userId)
	respondWithJSON(w, http.StatusOK, response{
		NewToken: newToken,
	})
}

package main

import (
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) postRevokeHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	refreshToken, found := strings.CutPrefix(authHeader, "Bearer ")

	if !found || refreshToken == "" {
		log.Printf("Did not find a valid authorization header")
		respondWithError(w, http.StatusUnauthorized, "Missing authorization")
		return
	}

	err := cfg.DB.RevokeRefreshToken(refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %q", err)
		respondWithError(w, http.StatusBadRequest, "Token does not exist or it already expired")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}

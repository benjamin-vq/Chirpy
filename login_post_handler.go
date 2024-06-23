package main

import (
	"encoding/json"
	"errors"
	"github.com/benjamin-vq/chirpy/internal/auth"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
)

type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

func (cfg *apiConfig) loginPostHandler(w http.ResponseWriter, r *http.Request) {

	loginParams := LoginParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loginParams)

	if err != nil {
		log.Printf("Error decoding login parameters: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode login parameters")
		return
	}

	user, err := cfg.DB.UserByEmail(loginParams.Email)

	if err != nil {
		if errors.Is(err, database.UserNotExists) {
			log.Printf("User with email %q does not exist %q", loginParams.Email, err)
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		log.Printf("Error finding user by email: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	err = auth.ComparePasswordHash(user.HashedPassword, loginParams.Password)
	if err != nil {
		log.Printf("Received password does not match")
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jwt, err := auth.CreateJwt(user.Id, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	rt, err := auth.GenerateRefreshToken()
	if err != nil {
		log.Printf("An error ocurred generating refresh token: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	err = cfg.DB.SaveToken(user.Id, rt)
	if err != nil {
		log.Printf("Error saving token to database after login: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not login")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		User: User{
			user.Email,
			user.Id,
		},
		Token:        jwt,
		RefreshToken: rt,
	})
}

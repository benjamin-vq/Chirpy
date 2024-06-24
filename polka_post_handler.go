package main

import (
	"encoding/json"
	"errors"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
	"strings"
)

type PolkaParams struct {
	Event string `json:"event"`
	Data  Data   `json:"data"`
}

type Data struct {
	UserId int `json:"user_id"`
}

func (cfg *apiConfig) postPolkaHandler(w http.ResponseWriter, r *http.Request) {

	authHeader := r.Header.Get("Authorization")
	token, found := strings.CutPrefix(authHeader, "ApiKey ")
	if !found || token != cfg.polkaApiKey {
		log.Printf("Invalid authorization header: %q", authHeader)
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	params := PolkaParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Could not decode polka params: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not decode params")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	user, err := cfg.DB.UserById(params.Data.UserId)
	if err != nil {
		log.Printf("Could not find user by id for polka webhook: %q", err)
		respondWithError(w, http.StatusNotFound, "")
		return
	}

	err = cfg.DB.MakeChirpyRed(user.Id)
	if err != nil {
		if errors.Is(err, database.UserNotExists) {
			log.Printf("Tried to update non-existing user to chirpy red: %q", err)
			respondWithError(w, http.StatusNotFound, "")
			return
		}
		log.Printf("Could not update user to chirpy red: %q", err)
		respondWithError(w, http.StatusInternalServerError, "Could not update user")
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
	return
}

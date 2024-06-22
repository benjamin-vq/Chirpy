package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		log.Printf("Error retrieving chirps from database: %q", err)
		respondWithError(w, 500, "Could not retrieve chirps")
		return
	}

	data, err := json.Marshal(chirps)

	if err != nil {
		log.Printf("Unable to unmarshal chirp list: %q", err)
		respondWithError(w, 500, "Could not to retrieve chirps")
		return
	}

	var httpStatus int
	if len(chirps) != 0 {
		httpStatus = http.StatusOK
	} else {
		httpStatus = http.StatusNoContent
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(data)
}

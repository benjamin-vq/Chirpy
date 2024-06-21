package main

import (
	"log"
	"net/http"
	"strconv"
)

func (cfg *apiConfig) chirpIdGetHandler(w http.ResponseWriter, r *http.Request) {

	p := r.PathValue("chirpId")
	id, err := strconv.Atoi(p)

	if err != nil {
		log.Printf("Unable to convert path value to a valid integer: %q", err)
		respondWithError(w, http.StatusBadRequest, "Provided id is not valid")
		return
	}

	chirp, err := cfg.DB.ChirpById(id)

	// TODO: Distinguish between 404 and 500
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}

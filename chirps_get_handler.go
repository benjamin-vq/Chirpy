package main

import (
	"cmp"
	"github.com/benjamin-vq/chirpy/internal/database"
	"log"
	"net/http"
	"slices"
	"strconv"
)

func (cfg *apiConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()

	if err != nil {
		log.Printf("Error retrieving chirps from database: %q", err)
		respondWithError(w, 500, "Could not retrieve chirps")
		return
	}

	idParamString := r.URL.Query().Get("author_id")
	sortParamString := r.URL.Query().Get("sort")
	var authorIdParam int
	if idParamString != "" {
		authorIdParam, err = strconv.Atoi(idParamString)
		if err != nil {
			log.Printf("Received an invalid author id as query param: %s", idParamString)
			respondWithError(w, http.StatusBadRequest, "Invalid author id")
			return
		}
	}

	filteredChirps := make([]database.Chirp, 0)
	for _, chirp := range chirps {
		// Continue to append chirps that don't need to be filtered
		if authorIdParam != 0 && chirp.AuthorId != authorIdParam {
			continue
		}
		filteredChirps = append(filteredChirps, chirp)
	}

	var httpStatus int
	if len(filteredChirps) != 0 {
		httpStatus = http.StatusOK
	} else {
		httpStatus = http.StatusNoContent
	}

	if sortParamString == "desc" {
		slices.SortFunc(filteredChirps, func(a, b database.Chirp) int {
			return cmp.Compare(b.Id, a.Id)
		})
	} else {
		slices.SortFunc(filteredChirps, func(a, b database.Chirp) int {
			return cmp.Compare(a.Id, b.Id)
		})
	}

	respondWithJSON(w, httpStatus, filteredChirps)
}

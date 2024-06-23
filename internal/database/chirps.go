package database

import (
	"errors"
	"fmt"
	"github.com/benjamin-vq/chirpy/internal/assert"
	"log"
)

type Chirp struct {
	Body     string `json:"body"`
	Id       int    `json:"id"`
	AuthorId int    `json:"author_id"`
}

var ChirpNotExists = errors.New("chirp does not exist")
var IncorrectAuthorId = errors.New("user id does not match chirp author id")

func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {

	//assert.That(body != "", "Chirp body can not be empty")
	assert.That(authorId != 0, "Should provide a valid author id")

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	// Ugly, but works. A better alternative would be to use another data structure for chirps
	var chirpId int
	for k, _ := range dbStructure.Chirps {
		if k > chirpId {
			chirpId = k
		}
	}
	// We just assigned to the latest id, increment.
	chirpId += 1
	chirp := Chirp{
		Body:     body,
		Id:       chirpId,
		AuthorId: authorId,
	}
	assert.That(dbStructure.Chirps != nil, "Chirps map should be initialized")
	dbStructure.Chirps[chirpId] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Printf("Error writing database structure: %q", err)
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database file to retrieve chirps: %q", err)
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, v := range dbStructure.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil
}

func (db *DB) ChirpById(id int) (Chirp, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database file to retrieve chirps: %q", err)
		return Chirp{}, err
	}

	chirp, exists := dbStructure.Chirps[id]

	if !exists {
		log.Printf("Chirp with id %d does not exist in database", id)
		return Chirp{}, fmt.Errorf("chirp with id %d does not exist", id)
	}

	return chirp, nil
}

func (db *DB) DeleteChirpById(chirpId, userId int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database file to delete chirp: %q", err)
		return err
	}

	chirp, exists := dbStructure.Chirps[chirpId]
	if !exists {
		log.Printf("Could not delete chirp with id %d because it does not exist", chirpId)
		return ChirpNotExists
	}

	if chirp.AuthorId != userId {
		log.Printf("Chirp author id (%d) does not match user id (%d)", chirp.AuthorId, userId)
		return IncorrectAuthorId
	}

	delete(dbStructure.Chirps, chirpId)
	log.Printf("Deleted chirp with author id %d from database", userId)

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Printf("Could not write database file to delete chirp: %q", err)
		return err
	}

	return nil
}

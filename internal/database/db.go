package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/benjamin-vq/chirpy/internal/assert"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type Chirp struct {
	Body string `json:"body"`
	Id   int    `json:"id"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {

	assert.That(path != "", "Database path can not be empty")

	db := DB{
		path: path,
		mu:   &sync.RWMutex{},
	}

	err := db.ensureDB()

	return &db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	//assert.That(body != "", "Chirp body can not be empty")

	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirpId := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		Body: body,
		Id:   chirpId,
	}
	assert.That(dbStructure.Chirps != nil, "Chirp map should be initialized")
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

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Database file does not exist, ensuring it exists by creating it")
		dbStructure := DBStructure{
			Chirps: make(map[int]Chirp),
		}
		db.writeDB(dbStructure)
		return nil
	}

	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	data, err := os.ReadFile(db.path)

	if err != nil {
		log.Printf("Could not read database file: %q", err)
		return DBStructure{}, err
	}

	dbStructure := DBStructure{}
	err = json.Unmarshal(data, &dbStructure)

	if err != nil {
		log.Printf("Could not unmarshal database structure: %q", err)
		return DBStructure{}, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(&dbStructure)

	if err != nil {
		log.Printf("Could not marshal database structure: %q", err)
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		log.Printf("Could not write structure to database file: %q", err)
		return err
	}

	log.Print("Succesfully wrote database structure to file")
	return nil

}

package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"

	"github.com/benjamin-vq/chirpy/internal/assert"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
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

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("Database file does not exist, ensuring it exists by creating it")
		dbStructure := DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
		err := db.writeDB(dbStructure)
		assert.That(err == nil, "Database could not be initialized: %q", err)
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

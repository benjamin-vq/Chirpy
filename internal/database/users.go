package database

import (
	"errors"
	"fmt"
	"github.com/benjamin-vq/chirpy/internal/assert"
	"log"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Id       int    `json:"id"`
}

var ErrEmailExists = fmt.Errorf("email already exists")

func (db *DB) CreateUser(email, password string) (User, error) {

	assert.That(email != "", "email can not be empty")

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database file to retrieve users: %q", err)
		return User{}, err
	}

	for id, user := range dbStructure.Users {
		if user.Email == email {
			log.Printf("Email %q already exists for user with id %d", email, id)
			return User{}, ErrEmailExists
		}
	}

	userId := len(dbStructure.Users) + 1
	user := User{
		email,
		password,
		userId,
	}

	assert.That(dbStructure.Users != nil, "Users map should be initialized")
	dbStructure.Users[userId] = user
	err = db.writeDB(dbStructure)

	if err != nil {
		log.Printf("Could not write database to save new user: %q", err)
		return User{}, err
	}

	log.Printf("Succesfully created user with id %d to database", user.Id)
	return user, nil
}

func (db *DB) UserByEmail(email string) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database file to retrieve users: %q", err)
		return User{}, err
	}

	for _, user := range dbStructure.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user does not exist")
}

package database

import (
	"errors"
	"github.com/benjamin-vq/chirpy/internal/assert"
	"log"
)

type User struct {
	Email          string `json:"email"`
	HashedPassword string `json:"hashedPassword"`
	Id             int    `json:"id"`
}

var ErrEmailExists = errors.New("email already exists")
var UserNotExists = errors.New("user does not exist")

func (db *DB) CreateUser(email, hashedPassword string) (User, error) {

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
		hashedPassword,
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

	return User{}, UserNotExists
}

func (db *DB) UserById(id int) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database to find a user by id: %q", err)
	}

	user, exists := dbStructure.Users[id]
	if !exists {
		log.Printf("User with id %d does not exist", id)
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(user *User) error {
	assert.That(user != nil, "Attempting to update nil user")

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database to find a user by id: %q", err)
		return err
	}

	if _, exists := dbStructure.Users[user.Id]; !exists {
		return UserNotExists
	}

	dbStructure.Users[user.Id] = *user
	err = db.writeDB(dbStructure)
	if err != nil {
		return err
	}

	log.Printf("Succesfully updated user in database")
	return nil
}

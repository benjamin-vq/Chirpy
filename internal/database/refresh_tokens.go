package database

import (
	"errors"
	"log"
	"time"
)

type RefreshToken struct {
	UserId    int       `json:"user_id"`
	Token     string    `json:"refresh_token"`
	ExpiresAt time.Time `json:"refresh_expires_at"`
}

func (db *DB) UserIdFromRefreshToken(rt string) (userId int, err error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database to save tokens: %q", err)
		return 0, err
	}

	refreshToken, exists := dbStructure.RefreshTokens[rt]
	if !exists {
		log.Print("Received token was not present in the database")
		return 0, errors.New("token does not exist")
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return 0, errors.New("refresh token expired")
	}

	return refreshToken.UserId, nil
}

func (db *DB) SaveToken(userId int, rt string) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database to save tokens: %q", err)
		return err
	}

	dbStructure.RefreshTokens[rt] = RefreshToken{
		UserId:    userId,
		Token:     rt,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Printf("Could not write database to save refresh token: %q", err)
		return err
	}

	return nil
}

func (db *DB) RevokeRefreshToken(rt string) error {

	dbStructure, err := db.loadDB()
	if err != nil {
		log.Printf("Could not load database to refresh token: %q", err)
		return err
	}

	if _, exists := dbStructure.RefreshTokens[rt]; exists {
		delete(dbStructure.RefreshTokens, rt)
		err = db.writeDB(dbStructure)
		if err != nil {
			log.Printf("Could not write database after revoking refresh token: %q", err)
			return err
		}
		return nil
	}

	return errors.New("did not find refresh token to revoke")
}

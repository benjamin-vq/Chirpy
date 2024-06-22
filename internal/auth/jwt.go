package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"log"
	"os"
	"strconv"
	"time"
)

func CreateJwt(expiration time.Duration, userId int) (string, error) {

	if expiration == 0 || expiration > 24*time.Hour {
		log.Printf("Jwt creation received an invalid expiration: %d. Defaulting to 24 hours", expiration)
		expiration = 24 * time.Hour
	}

	now := time.Now()

	issuedAt := jwt.NewNumericDate(now)
	expiresAt := jwt.NewNumericDate(issuedAt.Add(expiration))

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := os.Getenv("JWT_SECRET")
	signedJwt, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		log.Printf("Could not create signed JWT: %q", err)
		return "", err
	}
	log.Printf("Issued a new token at %v. Expires at %v", issuedAt, expiresAt)

	return signedJwt, nil
}

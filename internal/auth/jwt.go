package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

func CreateJwt(expiration time.Duration, userId int, jwtSecret string) (string, error) {

	if expiration == 0 || expiration > 24*time.Hour {
		log.Printf("Jwt creation received an invalid expiration: %d. Defaulting to 24 hours", expiration)
		expiration = 24 * time.Hour
	}

	now := time.Now().UTC()

	issuedAt := jwt.NewNumericDate(now)
	expiresAt := jwt.NewNumericDate(now.Add(expiration))

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedJwt, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		log.Printf("Could not create signed JWT: %q", err)
		return "", err
	}
	log.Printf("Issued a new token at %v. Expires at %v", issuedAt, expiresAt)

	return signedJwt, nil
}

func ValidateToken(token, jwtSecret string) (string, error) {

	claims := jwt.RegisteredClaims{}
	parsedJwt, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		log.Printf("Could not parse token: %q", err)
		return "", err
	}

	if !parsedJwt.Valid {
		log.Printf("Parsed JWT is invalid")
		return "", errors.New("invalid token")
	}

	subject, err := parsedJwt.Claims.GetSubject()
	if err != nil {
		log.Printf("Could not get subject in parsed token: %q", err)
		return "", err
	}

	issuer, err := parsedJwt.Claims.GetIssuer()
	if err != nil {
		log.Printf("Could not get issuer: %q", err)
		return "", err
	}
	if issuer != "chirpy" {
		log.Printf("Invalid issuer: %s", issuer)
		return "", nil
	}

	return subject, nil
}

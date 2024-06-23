package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"strconv"
	"time"
)

func CreateJwt(userId int, jwtSecret string) (string, error) {

	const expireAfter = 1 * time.Hour
	now := time.Now().UTC()

	issuedAt := jwt.NewNumericDate(now)
	expiresAt := jwt.NewNumericDate(now.Add(expireAfter))

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		Subject:   strconv.Itoa(userId),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	log.Printf("Issued a new token at %v. Expires at %v", issuedAt, expiresAt)

	return token.SignedString([]byte(jwtSecret))
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

func GenerateRefreshToken() (string, error) {

	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Printf("Could not generate refresh token: %q", err)
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

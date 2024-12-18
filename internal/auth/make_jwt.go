package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{Issuer: "Chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		Subject:   userID.String()}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := t.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return secret, nil
}

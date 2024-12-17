package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytePass := []byte(password)
	encryptedVal, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	res := string(encryptedVal)
	return res, nil
}

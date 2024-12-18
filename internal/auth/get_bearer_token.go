package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		err := errors.New("Authorization header missing or improperly formatted")
		return "", err
	}
	var trimBearer string
	if strings.HasPrefix(auth, "Bearer") {
		trimBearer = strings.TrimPrefix(auth, "Bearer ")
	} else {
		return "", errors.New("Bearer prefix not found or improperly formatted")
	}
	res := strings.TrimSpace(trimBearer)
	return res, nil
}

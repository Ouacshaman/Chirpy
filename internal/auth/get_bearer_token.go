package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		err := errors.New("Authorization not found in header")
		return "", err
	}
	trimBearer := strings.TrimPrefix(auth, "Bearer ")
	res := strings.TrimSpace(trimBearer)
	return res, nil
}

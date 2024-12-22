package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		err := errors.New("Authorization header missing or improperly formatted")
		return "", err
	}
	var apiKey string
	if strings.HasPrefix(auth, "ApiKey") {
		apiKey = strings.TrimPrefix(auth, "ApiKey ")
	} else {
		return "", errors.New("ApiKey prefix not found or improperly formatted")
	}
	res := strings.TrimSpace(apiKey)
	return res, nil
}

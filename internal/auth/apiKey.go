package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {

	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included")
	}
	splithAuth := strings.Split(authHeader, " ")
	if len(splithAuth) < 2 || splithAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splithAuth[1], nil
}

package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Example: Authorization: ApiKey {api code here}
func GetApiKey(headers http.Header) (string, error) {

	apiKey := headers.Get("Authorization")

	if apiKey == "" {
		return "", errors.New("não foram encontradas informações de autenticação")
	}

	auth := strings.Split(apiKey, " ")

	if len(auth) != 2 {
		return "", errors.New("header de autenticação contém problemas")
	}

	if auth[0] != "ApiKey" {
		return "", errors.New("primeira parte do Header de autenticação contém problemas")
	}

	return auth[1], nil
}

package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	apiKey, found := strings.CutPrefix(authHeader, "ApiKey")
	apiKey = strings.TrimSpace(apiKey)
	if !found || len(apiKey) == 0 {
		return "", fmt.Errorf("Authorization hearder has no API Key")
	}
	return apiKey, nil
}

package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	// Authorization: Bearer <token>
	authHeader := headers.Get("Authorization")
	tokenString, found := strings.CutPrefix(authHeader, "Bearer")
	if !found || len(tokenString) == 0 {
		return "", fmt.Errorf("Authorization hearder has no bearer token")
	}
	return strings.TrimSpace(tokenString), nil
}

package main

import (
	"net/http"
	"testing"

	"github.com/mshagirov/chirpy/internal/auth"
)

func TestBearerToken(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Errorf("Error creating a new http request: %v", err)
	}
	bearerToken := "token_string_goes_here"
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	got, err := auth.GetBearerToken(req.Header)
	if got != bearerToken || err != nil {
		t.Errorf(`Could not retrieve the bearer token, got
  %s
want
	%s`, got, bearerToken)
	}
	req.Header.Set("Authorization", "Bearer")
	if _, err := auth.GetBearerToken(req.Header); err == nil {
		t.Errorf("Expected error got nil")
	}
}

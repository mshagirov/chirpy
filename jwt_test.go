package main

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/chirpy/internal/auth"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "My secret message"
	tokenString, err := auth.MakeJWT(userID, secret, time.Second*30)

	if err != nil {
		t.Errorf("auth.MakeJWT(%v, %s,...) =\ntokenString : \"%s\"\nerr: %v\n",
			userID, secret, tokenString, err)
	}
	valUserID, err := auth.ValidateJWT(tokenString, secret)
	if err != nil {
		t.Errorf("auth.ValidateJWT(%s, %s) =\nUserID : \"%v\" , want \"%v\"\nerr: %v",
			tokenString, secret, valUserID, userID, err)
	}
	if userID != valUserID {
		t.Errorf("JWT: userID mismatch, got \"%v\" wanted \"%v\"", valUserID, userID)
	}
}

func TestJWTWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "My secret message"
	tokenString, err := auth.MakeJWT(userID, secret, time.Second*30)

	wrongSecret := "This is a wrong secret"
	someID, err := auth.ValidateJWT(tokenString, wrongSecret)
	if someID == userID {
		t.Errorf("Expected error got\nerr=%v\nwhen using a wrong tokenSecret\ngot\nsomeID=\"%v\"", err, someID)
	}
}

func TestJWTExpire(t *testing.T) {
	userID := uuid.New()
	secret := "My secret message"
	delay := time.Microsecond * 5
	tokenString, err := auth.MakeJWT(userID, secret, delay)
	time.Sleep(time.Millisecond * 100)
	someID, err := auth.ValidateJWT(tokenString, secret)
	if err == nil {
		t.Errorf("Expected error got\nerr=%v\nwhen waiting for token expiry tokenSecret\ngot\nsomeID=\"%v\"", err, someID)
	}
}

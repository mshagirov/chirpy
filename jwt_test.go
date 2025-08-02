package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mshagirov/chirpy/internal/auth"
)

func TestMakeJWT(t *testing.T) {
	fmt.Println("--- JWT create/validate ---")
	userID := uuid.New()
	fmt.Printf(`Created:
userID="%v"
`, userID)

	secret := "My secret message"
	fmt.Printf("tokenSecret=\"%s\"\n", secret)

	tokenString, err := auth.MakeJWT(userID, secret, time.Second*30)

	if err != nil {
		t.Errorf("auth.MakeJWT(%v, %s,...) =\ntokenString : \"%s\"\nerr: %v\n",
			userID, secret, tokenString, err)
	}
	fmt.Printf("Got tokenString=\n\"%s\"\n", tokenString)

	valUserID, err := auth.ValidateJWT(tokenString, secret)

	if err != nil {
		t.Errorf("auth.ValidateJWT(%s, %s) =\nUserID : \"%v\" , want \"%v\"\nerr: %v",
			tokenString, secret, valUserID, userID, err)
	}

	if userID != valUserID {
		t.Errorf("JWT: userID mismatch, got \"%v\" wanted \"%v\"", valUserID, userID)
	}
	fmt.Printf("Validation produced userID:\n%v\n", valUserID)
	fmt.Println("--- JWT create/validate ---")
}

func TestJWTWrongSecret(t *testing.T) {
	fmt.Println("--- JWT wrong secret ---")
	userID := uuid.New()
	fmt.Printf("Created:\nuserID=\"%v\"\n", userID)
	secret := "My secret message"
	fmt.Printf("tokenSecret=\"%s\"\n", secret)
	tokenString, err := auth.MakeJWT(userID, secret, time.Second*30)

	wrongSecret := "This is a wrong secret"
	fmt.Printf("Testing using a wrong secret: \"%s\"\n", wrongSecret)
	someID, err := auth.ValidateJWT(tokenString, wrongSecret)
	if someID == userID {
		t.Errorf("Expected error got\nerr=%v\nwhen using a wrong tokenSecret\ngot\nsomeID=\"%v\"", err, someID)
	}
	fmt.Printf("Got:\nuserID=%v, err=%v\n", someID, err)
	fmt.Println("--- JWT wrong secret ---")
}

func TestJWTExpire(t *testing.T) {
	fmt.Println("--- JWT expiry ---")
	userID := uuid.New()
	secret := "My secret message"
	delay := time.Microsecond * 5
	fmt.Println("ExpireIn:", delay.Microseconds())
	tokenString, err := auth.MakeJWT(userID, secret, delay)
	time.Sleep(time.Millisecond * 100)
	someID, err := auth.ValidateJWT(tokenString, secret)
	if err == nil {
		t.Errorf("Expected error got\nerr=%v\nwhen waiting for token expiry tokenSecret\ngot\nsomeID=\"%v\"", err, someID)
	}
	fmt.Printf("Got:\nuserID=%v, err=%v\n", someID, err)
}

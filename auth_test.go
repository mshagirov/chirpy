package main

import (
	"testing"

	"github.com/mshagirov/chirpy/internal/auth"
)

func TestHashPassword(t *testing.T) {
	s1 := "Hello world! 1, 2, 3"
	s2 := "hello world. 1, 2, 3"
	hash1, err := auth.HashPassword(s1)
	if err != nil {
		t.Errorf(`HashPassword(%s) = %s, %v`, s1, hash1, err)
	}
	hash1_duplicate, err := auth.HashPassword(s1)
	if err != nil || hash1 == hash1_duplicate {
		// different runs must not match
		t.Errorf(`
1) HashPassword(%s):
	%s
2) HashPassword(%s):
	%s
err=%v expect nil`, s1, hash1, s1, hash1_duplicate, err)
	}
	hash2, err := auth.HashPassword(s2)
	if err != nil || hash2 == hash1 {
		t.Errorf(`HashPassword:
1) HashPassword(%s):
	%s
2) HashPassword(%s):
	%s
err = %v expect nil`, s1, hash1, s2, hash2, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// CheckPasswordHash(password, hash string) error
	s := "Hello world! 1, 2, 3"
	hash, err := auth.HashPassword(s)
	if err != nil {
		t.Errorf(`HashPassword("%s") = "%s", %v`, s, hash, err)
	}
	if err := auth.CheckPasswordHash(s, hash); err != nil {
		t.Errorf(`CheckPasswordHash("%s", "%s") = %v`, s, hash, err)
	}
}

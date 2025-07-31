package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	s1 := "Hello world! 1, 2, 3"
	s2 := "hello world. 1, 2, 3"
	hash1, err := HashPassword(s1)
	if err != nil {
		t.Errorf(`HashPassword(%s) = %s, %v`, s1, hash1, err)
	}
	hash1_duplicate, err := HashPassword("Hello world!")
	if err != nil || hash1 != hash1_duplicat {
		t.Errorf(
			`Hash mismatch HashPassword(%s) = %s and %s, and/or err=%v`,
			s1, hash1, hash1_duplicate, err)
	}
	hash2, err := HashPassword(s2)
	if err != nil || hash2 == hash1 {
		t.Errorf(`HashPassword:
 %s hash:
			%s
 must be different from %s hash:
			%s
 or experiences an error:
			%v expected nil`,
			s1, hash1, s2, hash2, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// CheckPasswordHash(password, hash string) error
	s := "Hello world! 1, 2, 3"
	hash, err := HashPassword(s)
	if err != nil {
		t.Errorf(`HashPassword("%s") = "%s", %v`, s, hash, err)
	}
	if err := CheckPasswordHash(s, hash); err != nil {
		t.Errorf(`CheckPasswordHash("%s", "%s") = %v`, s, hash, err)
	}
}

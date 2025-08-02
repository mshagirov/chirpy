package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Issuer:    "chirpy",
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			Subject:   userID.String(),
		})
	jwtString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	var emptyUUID uuid.UUID
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims,
		func(token *jwt.Token) (any, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return emptyUUID, err
	}
	// userID
	userIDstr, err := token.Claims.GetSubject()
	if err != nil {
		return emptyUUID, err
	}
	userID, err := uuid.Parse(userIDstr)
	return userID, err
}

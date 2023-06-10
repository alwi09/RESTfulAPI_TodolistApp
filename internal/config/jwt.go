package config

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// untuk sign-in token
var JWTSecretKey = []byte("secret-key")

// payload token
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// create token
func CreateJWTToken(username string) (string, error) {
	// mengatur waktu kadalwarsa token
	expirationTime := time.Now().Add(10 * time.Minute)

	// membuat Claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// create token dengan sign-in method HS256 dan secret-key
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(JWTSecretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

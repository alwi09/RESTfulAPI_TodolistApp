package config

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// untuk sign-in token
const (
	jwtSecretKey       = "secret-key"
	jwtExpirationHours = 24
	jwtAlgorithm       = "HS256"
)

func CreateJWTToken(email string) (string, error) {
	// Buat payload token
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Hour * jwtExpirationHours).Unix(),
	}

	// Buat token JWT menggunakan kunci rahasia dan algoritma yang sesuai
	token := jwt.NewWithClaims(jwt.GetSigningMethod(jwtAlgorithm), claims)
	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

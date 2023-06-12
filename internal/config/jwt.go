package config

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// // untuk sign-in token
// var JWTSecretKey = []byte("secret-key")

// // payload token
// type Claims struct {
// 	Email string `json:"email"`
// 	jwt.StandardClaims
// }

// create token
// func CreateJWTToken(email string) (string, error) {
// 	// mengatur waktu kadalwarsa token
// 	expirationTime := time.Now().Add(24 * time.Hour)

// 	// membuat Claims
// 	claims := &Claims{
// 		Email: email,
// 		StandardClaims: jwt.StandardClaims{
// 			ExpiresAt: expirationTime.Unix(),
// 		},
// 	}

// 	// create token dengan sign-in method HS256 dan secret-key
// 	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

// 	// sign token dengan secret-key
// 	tokenString, err := token.SignedString(JWTSecretKey)

// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

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

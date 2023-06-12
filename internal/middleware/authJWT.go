package middleware

import (
	"net/http"
	"todolist_gin_gorm/internal/model/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// secret-key untuk sign-in token

// request -> server

// AuthMiddlewareJWT is a middleware function to check user authentication
func AuthMiddlewareJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// mengambil token dari Header Authorization
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "Authorization token not provided",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// split token dari Header
		tokenString := authHeader[len("Bearer "):]

		// Parsing token dengan menggunakan struct Claims
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret-key"), nil // Ganti dengan kunci rahasia yang sesuai
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
					Message: "unauthorized",
					Status:  http.StatusUnauthorized,
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{
				Message: "invalid or expired token",
				Status:  http.StatusBadRequest,
			})
			return
		}

		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "unauthorized (non valid)",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// jika token valid, mengambil email dari claim dan simpan ke dalam konteks
		ctx.Set("email", claims.Subject)

		// jika token valid, akan dilanjutkan ke handler
		ctx.Next()
	}
}

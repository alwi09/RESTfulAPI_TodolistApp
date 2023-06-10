package middleware

import (
	"net/http"
	"todolist_gin_gorm/internal/config"
	"todolist_gin_gorm/internal/model/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// secret-ket untuk sign-in token

// request -> server

func AuthMiddlewareJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// mengambil token dari Header Authorization
		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Message: "unauthorized",
				Status:  http.StatusUnauthorized,
			})
			return
		}

		// split token dari Header
		tokenString := authHeader[len("Bearer "):]

		// parsing token dengan secret-key
		claims := &config.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return config.JWTSecretKey, nil
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

		// jika token valid, mengambil username dari claim dan simpan ke dalam konteks
		ctx.Set("username", claims.Username)

		// jika token valid, akan di lanjutkan ke handler
		ctx.Next()
	}
}

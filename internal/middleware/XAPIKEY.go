package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func XAPIKEY() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader("Authorization")
		if key != "secret_lock" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{
					"message": "Unauthorized",
					"status":  http.StatusUnauthorized,
				})
			return
		}
	}
}

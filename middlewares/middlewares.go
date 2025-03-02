package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"jwt-psql/utils/token"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.ValidateToken(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

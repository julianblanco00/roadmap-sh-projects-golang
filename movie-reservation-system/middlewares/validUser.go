package middlewares

import (
	"movie-reservation-system/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInt := users.ExtractUserIdFromClaims(c)
		user := users.FindUserById(userIdInt)
		if user == nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

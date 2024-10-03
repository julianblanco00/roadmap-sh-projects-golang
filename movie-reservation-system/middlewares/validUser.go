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

func ValidAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdInt := users.ExtractUserIdFromClaims(c)
		error, role := users.ExtractRoleFromClaims(c)

		if error != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		user := users.FindUserById(userIdInt)
		if user == nil || user.Role != "admin" || role != "admin" {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

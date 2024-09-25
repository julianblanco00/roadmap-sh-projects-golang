package middlewares

import (
	"movie-reservation-system/users"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func ValidUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims := c.MustGet("user").(jwt.MapClaims)
		userIdInt, _ := strconv.Atoi(userClaims["_id"].(string))
		user := users.FindUserById(userIdInt)
		if user == nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

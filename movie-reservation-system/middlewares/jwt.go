package middlewares

import (
	"fmt"
	"movie-reservation-system/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := auth.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		fmt.Println(user.Claims)

		c.Next()
	}
}

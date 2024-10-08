package auth

import (
	"fmt"
	"movie-reservation-system/hashing"
	"movie-reservation-system/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleLogin(c *gin.Context) {
	email := c.Request.FormValue("email")
	password := c.Request.FormValue("password")

	user := users.FindUserByEmail(email)
	fmt.Println(user, email)

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !hashing.ComparePasswords(user.Password, password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})

		return
	}

	token, err := SignToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	welcomeMessage := "Welcome, " + user.Name
	c.JSON(http.StatusOK, gin.H{"message": welcomeMessage, "token": token})
}

package users

import (
	"fmt"
	"movie-reservation-system/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type User struct {
	ID        int
	Name      string
	Birthdate string
	Email     string
	Password  string
	Role      string
}

func ExtractUserIdFromClaims(c *gin.Context) int {
	userClaims := c.MustGet("user").(jwt.MapClaims)
	userIdInt, _ := strconv.Atoi(userClaims["_id"].(string))
	return userIdInt
}

func ExtractRoleFromClaims(c *gin.Context) (error, string) {
	userClaims := c.MustGet("user").(jwt.MapClaims)
	if userClaims["role"] == nil {
		return fmt.Errorf("Role not found"), ""
	}
	return nil, userClaims["role"].(string)
}

func FindUserByEmail(email string) *User {
	row := database.Db.QueryRow("SELECT * FROM users WHERE email = $1", email)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil
	}

	return &user
}

func FindUserById(id int) *User {
	row := database.Db.QueryRow("SELECT * FROM users WHERE id = $1", id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil
	}

	return &user
}

func GetUsers() ([]User, error) {
	rows, err := database.Db.Query("SELECT * FROM users")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Birthdate, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

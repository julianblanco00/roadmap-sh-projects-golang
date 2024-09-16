package database

import (
	"fmt"
	"movie-reservation-system/hashing"
)

type User struct {
	ID       int
	Name     string
	Age      int
	Email    string
	Password string
}

var users = []User{}

func CreateUser() (*User, error) {
	password, err := hashing.HashPassword("password")

	if err != nil {
		fmt.Println("Error hashing password")
		return nil, err
	}

	user := &User{
		Name:     "Test User",
		Age:      30,
		Email:    "testing1@gmail.com",
		Password: password,
	}

	users = append(users, *user)

	return user, nil
}

func FindUserByEmail(email string) *User {
	for _, user := range users {
		if user.Email == email {
			return &user
		}
	}

	return nil
}

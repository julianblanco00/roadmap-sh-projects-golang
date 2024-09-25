package main

import (
	"fmt"
	"log"
	"movie-reservation-system/auth"
	"movie-reservation-system/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func startWebServer() {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/auth/login", auth.HandleLogin)
	router.GET("/movies", database.GetMovies)

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	loadEnvVariables()
	database.Connect()
	startWebServer()
}

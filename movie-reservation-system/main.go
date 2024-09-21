package main

import (
	"fmt"
	"log"
	"movie-reservation-system/auth"
	"movie-reservation-system/database"
	"movie-reservation-system/middlewares"

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
	loadEnvVariables()
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.POST("/auth/login", auth.HandleLogin)
	router.GET("/movies", middlewares.JwtAuth(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "movies",
		})
		return
	})

	err := router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	database.Connect()
	startWebServer()
}

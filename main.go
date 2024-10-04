package main

import (
	"log"
	"login_page/database"
	"login_page/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	database.DbConnect()
}

func main() {
	// Set up Gin router
	r := gin.Default() // Initialize the Gin router here

	// Load HTML files
	r.LoadHTMLGlob("templates/*")

	routes.UserRoute(r)
	// routes.AdminRoute(r)

	// Start the server
	err := r.Run(":8089") // Corrected from R.run to r.Run
	if err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

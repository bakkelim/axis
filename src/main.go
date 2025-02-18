package main

import (
	"axis/src/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Set up middleware here if needed

	// Define routes
	routes.SetupRoutes(router)

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	router.Run(":" + port)
}

package routes

import (
	"axis/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Contract routes
		api.GET("/contracts/:id", controllers.GetContractByID)
		api.GET("/contracts/:id/execute", controllers.ExecuteContract)

		// Connector routes
		connectors := api.Group("/connectors")
		{
			connectors.POST("", controllers.CreateConnector)        // Create a new connector
			connectors.GET("", controllers.ListConnectors)          // List all connectors
			connectors.GET("/:id", controllers.GetConnector)        // Get a specific connector
			connectors.PUT("/:id", controllers.UpdateConnector)     // Update a connector
			connectors.DELETE("/:id", controllers.DeleteConnector)  // Delete a connector
			connectors.GET("/:id/test", controllers.TestConnection) // Test connection
		}
	}
}

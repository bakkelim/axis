package routes

import (
	"axis/src/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		// Contract routes
		contracts := api.Group("/contracts")
		{
			contracts.POST("", controllers.CreateContract)             // Create a new contract
			contracts.GET("", controllers.ListContracts)               // List all contracts
			contracts.GET("/:id", controllers.GetContractByID)         // Get a specific contract
			contracts.PUT("/:id", controllers.UpdateContract)          // Update a contract
			contracts.DELETE("/:id", controllers.DeleteContract)       // Delete a contract
			contracts.GET("/:id/execute", controllers.ExecuteContract) // Execute a contract
		}

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

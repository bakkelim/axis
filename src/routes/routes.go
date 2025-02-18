package routes

import (
	"github.com/gin-gonic/gin"
	"axis/src/controllers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.GET("/contracts/:id", controllers.GetContractByID)
		api.GET("/contracts/:id/execute", controllers.ExecuteContract)
	}
}

package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutes(t *testing.T) {
	router := gin.New()
	SetupRoutes(router)

	routesInfo := router.Routes()
	expectedRoutes := map[string]string{
		"/api/contracts/:id":         "GET",
		"/api/contracts/:id/execute": "GET",
	}

	for path, method := range expectedRoutes {
		found := false
		for _, route := range routesInfo {
			if route.Method == method && route.Path == path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected route %s %s not found", method, path)
		}
	}
}

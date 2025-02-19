package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)


func TestSetupRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	SetupRoutes(router)

	routesInfo := router.Routes()
	expectedRoutes := []struct {
		Method string
		Path   string
	}{
		// Contract routes
		{"GET", "/api/contracts/:id"},
		{"GET", "/api/contracts/:id/execute"},

		// Connector routes
		{"POST", "/api/connectors"},
		{"GET", "/api/connectors"},
		{"GET", "/api/connectors/:id"},
		{"PUT", "/api/connectors/:id"},
		{"DELETE", "/api/connectors/:id"},
		{"GET", "/api/connectors/:id/test"},
	}

	for _, expected := range expectedRoutes {
		found := false
		for _, route := range routesInfo {
			if route.Method == expected.Method && route.Path == expected.Path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected route %s %s not found", expected.Method, expected.Path)
		}
	}
}


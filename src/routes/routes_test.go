package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSetupRoutes(t *testing.T) {
	expectedRoutes := []struct {
		method string
		path   string
	}{
		{"GET", "/api/contracts"},
		{"POST", "/api/contracts"},
		{"GET", "/api/contracts/:id"},
		{"PUT", "/api/contracts/:id"},
		{"DELETE", "/api/contracts/:id"},
		{"POST", "/api/contracts/:id/execute"},
		{"POST", "/api/connectors"},
		{"GET", "/api/connectors"},
		{"GET", "/api/connectors/:id"},
		{"PUT", "/api/connectors/:id"},
		{"DELETE", "/api/connectors/:id"},
		{"GET", "/api/connectors/:id/test"},
	}

	// Create a new router
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Setup the routes
	SetupRoutes(router)

	// Check that all expected routes are registered
	routes := router.Routes()
	for _, expected := range expectedRoutes {
		found := false
		for _, route := range routes {
			if route.Method == expected.method && route.Path == expected.path {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected route %s %s not found", expected.method, expected.path)
		}
	}
}

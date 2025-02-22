package controllers

import (
	"axis/src/models"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// setupTestDirs creates test directories and returns their paths
func setupTestDirs(t *testing.T) (string, string, string) {
	baseDir := t.TempDir()
	contractsDir := filepath.Join(baseDir, "contracts")
	connectorsDir := filepath.Join(baseDir, "connectors")

	for _, dir := range []string{contractsDir, connectorsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
	}

	return baseDir, contractsDir, connectorsDir
}

// setupTestEnvironment creates temporary directories for both contracts and connectors
func setupTestEnvironment(t *testing.T) (string, string, string, func()) {
	baseDir, contractsDir, connectorsDir := setupTestDirs(t)

	// Set up test directories for both storage systems
	originalTestConnDir := testConnectorsDir
	testConnectorsDir = connectorsDir

	// Set up contract directory
	originalContractsDir := testContractsDir
	testContractsDir = contractsDir

	return baseDir, contractsDir, connectorsDir, func() {
		testConnectorsDir = originalTestConnDir
		testContractsDir = originalContractsDir
		os.RemoveAll(baseDir)
	}
}

// setupTestContext creates a test Gin context with the given request
func setupTestContext(method, path, body string) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	return c, w
}

// setupTestConnector creates a test connector and saves it
func setupTestConnector(t *testing.T, id string) *models.Connector {
	t.Logf("Setting up test connector with ID: %s", id)
	t.Logf("Current testConnectorsDir: %s", testConnectorsDir)

	connector := &models.Connector{
		ID:   id,
		Type: "postgres",
		Config: models.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
		},
	}

	// Verify the directory exists
	if _, err := os.Stat(testConnectorsDir); os.IsNotExist(err) {
		t.Fatalf("Test connectors directory does not exist: %s", testConnectorsDir)
	}

	// Save the connector
	t.Logf("Saving connector to directory: %s", testConnectorsDir)
	if err := saveConnector(connector); err != nil {
		t.Logf("Failed to save connector: %v", err)
		t.Fatalf("Failed to save connector: %v", err)
	}

	// Verify the file was created
	filePath := filepath.Join(testConnectorsDir, id+".json")
	t.Logf("Checking for connector file at: %s", filePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Fatalf("Connector file was not created: %s", filePath)
	}

	// Try to read the file back to verify its contents
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read connector file: %v", err)
	}
	t.Logf("Connector file contents: %s", string(data))

	return connector
}

package controllers

import (
	"axis/src/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// To allow testing, we assume saveConnector is a package level variable.
// We'll override it in tests and restore after each test.
var originalSaveConnector = saveConnector

func restoreSaveConnector() {
	saveConnector = originalSaveConnector
}

// setupTestConnectorRequest creates a test connector request with the given ID
func setupTestConnectorRequest(method, path string, connector *models.Connector) (*gin.Context, *httptest.ResponseRecorder) {
	data, _ := json.Marshal(connector)
	c, w := setupTestContext(method, path, string(data))
	// Extract ID from path, handling both /connectors/id and /connectors/id/test paths
	parts := strings.Split(path, "/")
	if len(parts) >= 3 && parts[2] != "" {
		// parts[0] is empty because path starts with /
		// parts[1] is "connectors"
		// parts[2] is the ID
		id := parts[2]
		c.Params = gin.Params{{Key: "id", Value: id}}
	}
	return c, w
}

func TestCreateConnector_Success(t *testing.T) {
	defer restoreSaveConnector()

	// Override saveConnector to simulate success.
	saveConnector = func(connector *models.Connector) error {
		return nil
	}

	connector := &models.Connector{
		Name: "Test Connector",
	}
	c, w := setupTestConnectorRequest("POST", "/connectors", connector)
	CreateConnector(c)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Connector
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, "Test Connector", response.Name)
}

func TestCreateConnector_InvalidJSON(t *testing.T) {
	c, w := setupTestContext("POST", "/connectors", "invalid json")
	CreateConnector(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateConnector_SaveError(t *testing.T) {
	defer restoreSaveConnector()

	// Override saveConnector to simulate a failure.
	saveConnector = func(connector *models.Connector) error {
		return errors.New("save error")
	}

	connector := &models.Connector{
		ID:   "test-conn",
		Name: "Test Connector",
	}
	c, w := setupTestConnectorRequest("POST", "/connectors", connector)
	CreateConnector(c)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateConnector(t *testing.T) {
	connector := &models.Connector{
		Type: "postgres",
		Config: models.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
		},
	}
	c, w := setupTestConnectorRequest("POST", "/connectors", connector)

	CreateConnector(c)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Connector
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.ID)
	assert.Equal(t, connector.Type, response.Type)
	assert.Equal(t, connector.Config, response.Config)
}

// Test GetConnector
func TestGetConnector(t *testing.T) {
	_, _, connectorsDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	t.Logf("Using connectors directory: %s", connectorsDir)
	t.Logf("testConnectorsDir is set to: %s", testConnectorsDir)

	connector := setupTestConnector(t, "test-conn")
	t.Logf("Connector saved: %+v", connector)

	c, w := setupTestConnectorRequest("GET", "/connectors/test-conn", nil)

	GetConnector(c)
	t.Logf("Response status: %d, body: %s", w.Code, w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Connector
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, connector.ID, response.ID)
}

func TestGetConnector_NotFound(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	c, w := setupTestConnectorRequest("GET", "/connectors/nonexistent", nil)
	GetConnector(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test ListConnectors
func TestListConnectors_Success(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	setupTestConnector(t, "conn1")
	setupTestConnector(t, "conn2")
	c, w := setupTestContext("GET", "/connectors", "")
	ListConnectors(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Connector
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
}

func TestListConnectors_Empty(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	c, w := setupTestContext("GET", "/connectors", "")
	ListConnectors(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.Connector
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Empty(t, response)
}

// Test UpdateConnector
func TestUpdateConnector(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	connector := setupTestConnector(t, "test-conn")
	connector.Type = "mysql" // Change something
	c, w := setupTestConnectorRequest("PUT", "/connectors/test-conn", connector)
	UpdateConnector(c)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateConnector_InvalidJSON(t *testing.T) {
	c, w := setupTestConnectorRequest("PUT", "/connectors/test-conn", nil)
	c.Request = httptest.NewRequest("PUT", "/connectors/test-conn", strings.NewReader("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")
	UpdateConnector(c)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateConnector_NotFound(t *testing.T) {
	connector := &models.Connector{ID: "nonexistent"}
	c, w := setupTestConnectorRequest("PUT", "/connectors/nonexistent", connector)
	UpdateConnector(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test DeleteConnector
func TestDeleteConnector_Success(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_ = setupTestConnector(t, "test-conn")
	c, w := setupTestConnectorRequest("DELETE", "/connectors/test-conn", nil)
	DeleteConnector(c)
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify connector was deleted
	c2, w2 := setupTestConnectorRequest("GET", "/connectors/test-conn", nil)
	GetConnector(c2)
	assert.Equal(t, http.StatusNotFound, w2.Code)
}

func TestDeleteConnector_NotFound(t *testing.T) {
	c, w := setupTestConnectorRequest("DELETE", "/connectors/nonexistent", nil)
	DeleteConnector(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Test TestConnection
func TestTestConnection(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	_ = setupTestConnector(t, "test-conn")
	c, w := setupTestConnectorRequest("GET", "/connectors/test-conn/test", nil)
	TestConnection(c)
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "connection successful", response["status"])
}

func TestTestConnection_NotFound(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	c, w := setupTestConnectorRequest("GET", "/connectors/nonexistent/test", nil)
	TestConnection(c)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

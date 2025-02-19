package controllers

import (
	"axis/src/models"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
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

func TestCreateConnector_Success(t *testing.T) {
	defer restoreSaveConnector()

	// Override saveConnector to simulate success.
	saveConnector = func(connector *models.Connector) error {
		return nil
	}

	// Set Gin in test mode.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/connectors", CreateConnector)

	// Create a valid JSON body; adjust fields as per your models.Connector definition.
	connectorData := map[string]interface{}{
		"name": "Test Connector",
	}
	body, _ := json.Marshal(connectorData)

	req, err := http.NewRequest(http.MethodPost, "/connectors", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.Connector
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	// Check that an ID was generated.
	assert.NotEmpty(t, response.ID)
	// Optionally check that the Name field is preserved.
	assert.Equal(t, "Test Connector", response.Name)
}

func TestCreateConnector_InvalidJSON(t *testing.T) {
	// No need to override saveConnector since the JSON is invalid.
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/connectors", CreateConnector)

	// Use an invalid JSON body.
	req, err := http.NewRequest(http.MethodPost, "/connectors", bytes.NewBufferString("invalid json"))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateConnector_SaveError(t *testing.T) {
	defer restoreSaveConnector()

	// Override saveConnector to simulate a failure.
	saveConnector = func(connector *models.Connector) error {
		return errors.New("save error")
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/connectors", CreateConnector)

	connectorData := map[string]interface{}{
		"name": "Test Connector",
	}
	body, _ := json.Marshal(connectorData)

	req, err := http.NewRequest(http.MethodPost, "/connectors", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

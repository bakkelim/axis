package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"axis/src/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestEnvironment(t *testing.T) (string, string, string) {
	baseDir := t.TempDir()
	dataContractsDir := filepath.Join(baseDir, "data-contracts")
	controllersDir := filepath.Join(baseDir, "controllers")

	require.NoError(t, os.MkdirAll(dataContractsDir, os.ModePerm))
	require.NoError(t, os.MkdirAll(controllersDir, os.ModePerm))

	return baseDir, dataContractsDir, controllersDir
}

func TestGetContractByID_Success(t *testing.T) {
	_, dataContractsDir, controllersDir := setupTestEnvironment(t)

	// Create sample contract
	contract := models.Contract{
		ID:   "123",
		Name: "Test Contract",
	}
	contractBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	// Write contract file
	contractPath := filepath.Join(dataContractsDir, "contract-123.json")
	require.NoError(t, os.WriteFile(contractPath, contractBytes, 0644))

	// Setup working directory
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(controllersDir))
	defer os.Chdir(origDir)

	// Setup Gin test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	// Execute test
	GetContractByID(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	var response models.Contract
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, contract.ID, response.ID)
	assert.Equal(t, contract.Name, response.Name)
}

func TestExecuteContract_DBConnectionFailed(t *testing.T) {
	_, dataContractsDir, controllersDir := setupTestEnvironment(t)

	// Create sample contract with invalid DB connector
	contract := models.Contract{
		ID: "456",
		Connector: models.Connector{
			Type:             "invalid",
			ConnectionString: "",
			SQLQuery:         "SELECT 1",
		},
		ResponseTemplate: models.ResponseTemplate{
			Template: map[string]interface{}{
				"result": "{{.result}}",
			},
		},
	}
	contractBytes, err := json.Marshal(contract)
	require.NoError(t, err)

	// Write contract file
	contractPath := filepath.Join(dataContractsDir, "contract-456.json")
	require.NoError(t, os.WriteFile(contractPath, contractBytes, 0644))

	// Setup working directory
	origDir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(controllersDir))
	defer os.Chdir(origDir)

	// Setup Gin test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "456"}}

	// Execute test
	ExecuteContract(c)

	// Verify response
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, "Database connection failed", response["error"])
}

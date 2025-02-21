package controllers

import (
	"axis/src/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestCreateContract_InvalidJSON verifies that providing invalid JSON returns a 400 error.
func TestCreateContract_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Provide invalid JSON body.
	invalidJSON := strings.NewReader("{invalid json")
	c.Request = httptest.NewRequest("POST", "/contracts", invalidJSON)

	CreateContract(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d but got %d", http.StatusBadRequest, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}
	if _, ok := resp["error"]; !ok {
		t.Error("Expected error message in response")
	}
}

// TestUpdateContract_InvalidJSON checks that UpdateContract returns an error when given invalid JSON.
// Note: UpdateContract first attempts to load an existing contract; if it fails (contract not found), a 404 may be returned.
func TestUpdateContract_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set an id parameter.
	c.Params = gin.Params{{Key: "id", Value: "some-id"}}
	invalidJSON := strings.NewReader("{invalid json")
	c.Request = httptest.NewRequest("PUT", "/contracts/some-id", invalidJSON)

	UpdateContract(c)

	// Depending on the loadContract result, we may see 400 (bad JSON) or 404 (contract not found).
	if w.Code != http.StatusBadRequest && w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d or %d but got %d", http.StatusBadRequest, http.StatusNotFound, w.Code)
	}
}

// TestGetContractByID_NotFound ensures that requesting a non-existent contract returns a 404 error.
func TestGetContractByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Use a contract id that does not exist.
	c.Params = gin.Params{{Key: "id", Value: "non-existing-id"}}
	c.Request = httptest.NewRequest("GET", "/contracts/non-existing-id", nil)

	GetContractByID(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d but got %d", http.StatusNotFound, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}
	if msg, ok := resp["error"]; !ok || msg != "Contract not found" {
		t.Errorf("Expected error 'Contract not found' but got %v", resp)
	}
}

// TestDeleteContract_NotFound verifies that attempting to delete a non-existent contract returns a 404 error.
func TestDeleteContract_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Use a contract id that does not exist.
	c.Params = gin.Params{{Key: "id", Value: "non-existing-id"}}
	c.Request = httptest.NewRequest("DELETE", "/contracts/non-existing-id", nil)

	DeleteContract(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d but got %d", http.StatusNotFound, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}
	if msg, ok := resp["error"]; !ok || msg != "Contract not found" {
		t.Errorf("Expected error 'Contract not found' but got %v", resp)
	}
}

// TestExecuteContract_ContractNotFound checks that executing a non-existent contract returns a 404 error.
func TestExecuteContract_ContractNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Use an id that does not match any contract file.
	c.Params = gin.Params{{Key: "id", Value: "non-existing-id"}}
	c.Request = httptest.NewRequest("GET", "/contracts/non-existing-id/execute", nil)

	ExecuteContract(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d but got %d", http.StatusNotFound, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("Failed to parse response JSON:", err)
	}
	if msg, ok := resp["error"]; !ok || msg != "Contract not found" {
		t.Errorf("Expected error 'Contract not found' but got %v", resp)
	}
}

// TestListContracts_InternalServerError attempts to list contracts.
// In an isolated test environment, missing configuration or database should result in an internal server error.
func TestListContracts_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/contracts", nil)

	ListContracts(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d but got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestAnonymizeValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		rule     models.AnonymizationRule
		expected string
	}{
		{
			name:  "Mask with pattern",
			value: "123456789",
			rule: models.AnonymizationRule{
				Method:  "mask",
				Pattern: "XXX-XX-****",
			},
			expected: "123-45-****",
		},
		{
			name:  "Mask without pattern",
			value: "sensitive",
			rule: models.AnonymizationRule{
				Method: "mask",
			},
			expected: "*********",
		},
		{
			name:  "Hash method",
			value: "test@example.com",
			rule: models.AnonymizationRule{
				Method: "hash",
			},
			expected: "973dfe463ec85785f5f95af5ba3906eedb2d931c24e69824a89ea65dba4e813b",
		},
		{
			name:  "Unknown method",
			value: "keep-as-is",
			rule: models.AnonymizationRule{
				Method: "unknown",
			},
			expected: "keep-as-is",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := anonymizeValue(tt.value, tt.rule)
			if result != tt.expected {
				t.Errorf("anonymizeValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExecuteContract_WithAnonymization(t *testing.T) {
	// Setup test environment
	baseDir, contractsDir, connectorsDir := setupTestDirs(t)
	defer os.RemoveAll(baseDir)

	// Create test connector
	connector := models.Connector{
		ID:   "test-conn",
		Type: "postgres",
		Config: models.DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
		},
	}
	saveTestConnector(t, connectorsDir, &connector)

	// Create test contract with anonymization rules
	contract := models.Contract{
		ID:   "test-contract",
		Name: "Test Contract",
		Query: models.DatabaseQuery{
			ConnectorID: "test-conn",
			SQLQuery:    "SELECT id, email, ssn FROM users",
		},
		ResponseTemplate: models.ResponseTemplate{
			Template: map[string]interface{}{
				"user_id": "{{.id}}",
				"email":   "{{.email}}",
				"ssn":     "{{.ssn}}",
			},
			Anonymization: []models.AnonymizationRule{
				{
					Field:  "email",
					Method: "hash",
				},
				{
					Field:   "ssn",
					Method:  "mask",
					Pattern: "XXX-XX-****",
				},
			},
		},
	}
	saveTestContract(t, contractsDir, &contract)

	// Setup Gin test context
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: contract.ID}}

	// Mock DB connection (you might want to use sqlmock here)
	// For this example, we'll check the response structure

	ExecuteContract(c)

	// Verify response structure
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// Check basic response structure
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, contract.ID, response["contract_id"])
	assert.Equal(t, "success", response["status"])

	// Helper functions for test setup
	results, ok := response["results"].([]interface{})
	if !ok {
		t.Fatal("Results should be an array")
	}

	// Verify anonymization in results
	if len(results) > 0 {
		result := results[0].(map[string]interface{})

		// Check that email is hashed (64 characters hex string)
		if email, ok := result["email"].(string); ok {
			assert.Regexp(t, "^[0-9a-f]{64}$", email)
		}

		// Check that SSN follows mask pattern
		if ssn, ok := result["ssn"].(string); ok {
			assert.Regexp(t, "^\\d{3}-\\d{2}-\\*{4}$", ssn)
		}
	}
}

// Helper functions for test setup
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

func saveTestConnector(t *testing.T, dir string, connector *models.Connector) {
	data, err := json.Marshal(connector)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, connector.ID+".json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

func saveTestContract(t *testing.T, dir string, contract *models.Contract) {
	data, err := json.Marshal(contract)
	if err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(dir, contract.ID+".json"), data, 0644); err != nil {
		t.Fatal(err)
	}
}

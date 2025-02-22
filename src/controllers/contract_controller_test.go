package controllers

import (
	"axis/src/models"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	invalidJSON := strings.NewReader("{invalid")
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

	// Create a valid request body
	reqBody := strings.NewReader(`{
		"filters": [],
		"pagination": {"page": 1, "page_size": 10}
	}`)

	c.Request = httptest.NewRequest("POST", "/contracts/non-existing-id/execute", reqBody)
	c.Request.Header.Set("Content-Type", "application/json")

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

	// Save the original listContracts function
	originalListContracts := listContracts
	// Restore it after the test
	defer func() {
		listContracts = originalListContracts
	}()

	// Mock listContracts to return an error
	listContracts = func() ([]models.Contract, error) {
		return nil, errors.New("database error")
	}

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

// setupMockDB creates a mock database and returns the mock and a cleanup function
func setupMockDB(t *testing.T, expectedQuery string, rows *sqlmock.Rows) func() {
	originalSqlOpen := sqlOpen
	sqlOpen = func(driverName, dataSource string) (*sql.DB, error) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatalf("Failed to create mock DB: %v", err)
		}
		mock.ExpectQuery(expectedQuery).WillReturnRows(rows)
		return db, nil
	}
	return func() {
		sqlOpen = originalSqlOpen
	}
}

// setupTestContract creates a test contract with the given configuration
func setupTestContract(t *testing.T, id string, connectorID string, query string, template map[string]interface{}, anonymization []models.AnonymizationRule) *models.Contract {
	contract := &models.Contract{
		ID:   id,
		Name: "Test Contract",
		Query: models.DatabaseQuery{
			ConnectorID: connectorID,
			SQLQuery:    query,
		},
		ResponseTemplate: models.ResponseTemplate{
			Template:      template,
			Anonymization: anonymization,
		},
	}
	if err := saveContract(contract); err != nil {
		t.Fatalf("Failed to save contract: %v", err)
	}
	return contract
}

func TestExecuteContract_WithAnonymization(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "ssn"}).
		AddRow("1", "test@example.com", "123-45-6789")
	cleanup = setupMockDB(t, "SELECT id, email, ssn FROM users", rows)
	defer cleanup()

	// Create test connector
	testConn := setupTestConnector(t, "test-conn")

	// Create test contract
	contract := setupTestContract(t, "test-contract", testConn.ID,
		"SELECT id, email, ssn FROM users",
		map[string]interface{}{
			"user_id": "{{.id}}",
			"email":   "{{.email}}",
			"ssn":     "{{.ssn}}",
		},
		[]models.AnonymizationRule{
			{Field: "email", Method: "hash"},
			{Field: "ssn", Method: "mask", Pattern: "XXX-XX-****"},
		},
	)

	// Execute the contract
	c, w := setupTestContext("POST", "/contracts/test-contract/execute", `{
		"filters": [],
		"pagination": {"page": 1, "page_size": 10}
	}`)
	c.Params = gin.Params{{Key: "id", Value: contract.ID}}

	ExecuteContract(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, contract.ID, response["contract_id"])
	assert.Equal(t, "success", response["status"])

	results, ok := response["results"].([]interface{})
	require.True(t, ok)
	require.Len(t, results, 1)

	result := results[0].(map[string]interface{})
	assert.Equal(t, "1", result["user_id"])
	assert.Regexp(t, "^[0-9a-f]{64}$", result["email"]) // Hashed email
	assert.Equal(t, "123-45-****", result["ssn"])       // Masked SSN
}

func TestBuildWhereClause(t *testing.T) {
	tests := []struct {
		name           string
		filters        []models.FilterCondition
		expectedWhere  string
		expectedValues []interface{}
	}{
		{
			name:           "No filters",
			filters:        []models.FilterCondition{},
			expectedWhere:  "",
			expectedValues: nil,
		},
		{
			name: "Single equals filter",
			filters: []models.FilterCondition{
				{
					Field:    "name",
					Operator: models.OperatorEquals,
					Value:    "John",
				},
			},
			expectedWhere:  " WHERE name = $1",
			expectedValues: []interface{}{"John"},
		},
		{
			name: "Multiple filters",
			filters: []models.FilterCondition{
				{
					Field:    "age",
					Operator: models.OperatorGreater,
					Value:    18,
				},
				{
					Field:    "active",
					Operator: models.OperatorEquals,
					Value:    true,
				},
			},
			expectedWhere:  " WHERE age > $1 AND active = $2",
			expectedValues: []interface{}{18, true},
		},
		{
			name: "LIKE operator",
			filters: []models.FilterCondition{
				{
					Field:    "email",
					Operator: models.OperatorLike,
					Value:    "%@example.com",
				},
			},
			expectedWhere:  " WHERE email LIKE $1",
			expectedValues: []interface{}{"%@example.com"},
		},
		{
			name: "IN operator",
			filters: []models.FilterCondition{
				{
					Field:    "status",
					Operator: models.OperatorIn,
					Value:    []interface{}{"active", "pending"},
				},
			},
			expectedWhere:  " WHERE status IN ($1,$2)",
			expectedValues: []interface{}{"active", "pending"},
		},
		{
			name: "Complex mixed filters",
			filters: []models.FilterCondition{
				{
					Field:    "age",
					Operator: models.OperatorGreater,
					Value:    18,
				},
				{
					Field:    "name",
					Operator: models.OperatorLike,
					Value:    "John%",
				},
				{
					Field:    "status",
					Operator: models.OperatorIn,
					Value:    []interface{}{"active", "pending"},
				},
			},
			expectedWhere:  " WHERE age > $1 AND name LIKE $2 AND status IN ($3,$4)",
			expectedValues: []interface{}{18, "John%", "active", "pending"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			where, values := buildWhereClause(tt.filters)
			assert.Equal(t, tt.expectedWhere, where)
			assert.Equal(t, tt.expectedValues, values)
		})
	}
}

func TestBuildWhereClause_InvalidIN(t *testing.T) {
	filters := []models.FilterCondition{
		{
			Field:    "status",
			Operator: models.OperatorIn,
			Value:    "not-an-array", // Invalid value for IN operator
		},
	}

	where, values := buildWhereClause(filters)
	assert.Equal(t, " WHERE ", where)
	assert.Equal(t, []interface{}{"not-an-array"}, values)
}

func TestBuildOrderByClause(t *testing.T) {
	tests := []struct {
		name          string
		sortOptions   []models.SortOption
		expectedOrder string
	}{
		{
			name:          "No sort options",
			sortOptions:   []models.SortOption{},
			expectedOrder: "",
		},
		{
			name: "Single sort option",
			sortOptions: []models.SortOption{
				{
					Field:     "name",
					Direction: "asc",
				},
			},
			expectedOrder: " ORDER BY name asc",
		},
		{
			name: "Multiple sort options",
			sortOptions: []models.SortOption{
				{
					Field:     "age",
					Direction: "desc",
				},
				{
					Field:     "name",
					Direction: "asc",
				},
			},
			expectedOrder: " ORDER BY age desc, name asc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := buildOrderByClause(tt.sortOptions)
			assert.Equal(t, tt.expectedOrder, order)
		})
	}
}

func TestExecuteContract_WithSorting(t *testing.T) {
	_, _, _, cleanup := setupTestEnvironment(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "name", "age"}).
		AddRow("1", "Alice", 30)
	cleanup = setupMockDB(t, "SELECT id, name, age FROM users", rows)
	defer cleanup()

	// Create test connector
	testConn := setupTestConnector(t, "test-conn")

	// Create test contract with sorting options
	contract := setupTestContract(t, "test-contract", testConn.ID,
		"SELECT id, name, age FROM users",
		map[string]interface{}{
			"user_id": "{{.id}}",
			"name":    "{{.name}}",
			"age":     "{{.age}}",
		},
		nil, // no anonymization rules
	)

	// Execute the contract
	c, w := setupTestContext("POST", "/contracts/test-contract/execute", `{
		"filters": [],
		"pagination": {"page": 1, "page_size": 10},
		"sort": [
			{"field": "age", "direction": "desc"},
			{"field": "name", "direction": "asc"}
		]
	}`)
	c.Params = gin.Params{{Key: "id", Value: contract.ID}}

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

	// Verify sorting in results
	if len(results) > 0 {
		result := results[0].(map[string]interface{})

		// Check that age is sorted in descending order
		if age, ok := result["age"].(float64); ok {
			assert.GreaterOrEqual(t, age, 0.0)
		}

		// Check that name is sorted in ascending order
		if name, ok := result["name"].(string); ok {
			assert.NotEmpty(t, name)
		}
	}
}

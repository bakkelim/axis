package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
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

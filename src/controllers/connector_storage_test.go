package controllers

import (
	"axis/src/models"
	"testing"
)

// setupTestConnectorDir creates a temporary directory for connector tests and returns a cleanup function
func setupTestConnectorDir(t *testing.T) func() {
	tmpDir := t.TempDir()
	originalTestDir := testConnectorsDir
	testConnectorsDir = tmpDir
	return func() {
		testConnectorsDir = originalTestDir
	}
}

func TestSaveAndLoadConnector(t *testing.T) {
	cleanup := setupTestConnectorDir(t)
	defer cleanup()

	// Create test connector
	connector := &models.Connector{
		ID: "test1",
		// Add other required fields...
	}

	// Test saving
	if err := saveConnector(connector); err != nil {
		t.Fatalf("saveConnector failed: %v", err)
	}

	// Test loading
	loaded, err := LoadConnector(connector.ID)
	if err != nil {
		t.Fatalf("loadConnector failed: %v", err)
	}

	if loaded.ID != connector.ID {
		t.Errorf("Expected ID %s, got %s", connector.ID, loaded.ID)
	}
}

func TestListConnectors(t *testing.T) {
	cleanup := setupTestConnectorDir(t)
	defer cleanup()

	connector1 := &models.Connector{ID: "list1"}
	connector2 := &models.Connector{ID: "list2"}

	if err := saveConnector(connector1); err != nil {
		t.Fatalf("saveConnector for connector1 failed: %v", err)
	}
	if err := saveConnector(connector2); err != nil {
		t.Fatalf("saveConnector for connector2 failed: %v", err)
	}

	connectors, err := listConnectors()
	if err != nil {
		t.Fatalf("listConnectors failed: %v", err)
	}

	found1, found2 := false, false
	for _, c := range connectors {
		if c.ID == "list1" {
			found1 = true
		}
		if c.ID == "list2" {
			found2 = true
		}
	}
	if !found1 || !found2 {
		t.Fatalf("Expected to find connectors with IDs 'list1' and 'list2', got: %#v", connectors)
	}
}

func TestDeleteConnector(t *testing.T) {
	cleanup := setupTestConnectorDir(t)
	defer cleanup()

	connector := &models.Connector{ID: "delete1"}

	if err := saveConnector(connector); err != nil {
		t.Fatalf("saveConnector failed: %v", err)
	}

	if err := deleteConnector("delete1"); err != nil {
		t.Fatalf("deleteConnector failed: %v", err)
	}

	_, err := LoadConnector("delete1")
	if err == nil {
		t.Fatalf("Expected error when loading a deleted connector, got nil")
	}
	if err.Error() != "connector not found" {
		t.Fatalf("Expected error 'connector not found', got: %v", err)
	}
}

func TestLoadNonexistentConnector(t *testing.T) {
	cleanup := setupTestConnectorDir(t)
	defer cleanup()

	_, err := LoadConnector("nonexistent")
	if err == nil {
		t.Fatalf("Expected error for nonexistent connector, got nil")
	}
	if err.Error() != "connector not found" {
		t.Fatalf("Expected error 'connector not found', got: %v", err)
	}
}

package controllers

import (
	"axis/src/models"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestMain ensures a clean environment for our tests.
func TestMain(m *testing.M) {
	// connectorsDir is "../connectors" relative to this file.
	dir := filepath.Join("..", "connectors")
	_ = os.RemoveAll(dir)
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func TestSaveAndLoadConnector(t *testing.T) {
	connector := &models.Connector{
		ID: "test1",
		// add other fields as defined in your models.Connector if needed
	}

	if err := saveConnector(connector); err != nil {
		t.Fatalf("saveConnector failed: %v", err)
	}

	loaded, err := loadConnector("test1")
	if err != nil {
		t.Fatalf("loadConnector failed: %v", err)
	}

	if !reflect.DeepEqual(connector, loaded) {
		t.Fatalf("loaded connector does not match saved connector.\nExpected: %#v\nGot: %#v", connector, loaded)
	}
}

func TestListConnectors(t *testing.T) {
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
	connector := &models.Connector{ID: "delete1"}

	if err := saveConnector(connector); err != nil {
		t.Fatalf("saveConnector failed: %v", err)
	}

	if err := deleteConnector("delete1"); err != nil {
		t.Fatalf("deleteConnector failed: %v", err)
	}

	_, err := loadConnector("delete1")
	if err == nil {
		t.Fatalf("Expected error when loading a deleted connector, got nil")
	}
	if err.Error() != "connector not found" {
		t.Fatalf("Expected error 'connector not found', got: %v", err)
	}
}

func TestLoadNonexistentConnector(t *testing.T) {
	_, err := loadConnector("nonexistent")
	if err == nil {
		t.Fatalf("Expected error for nonexistent connector, got nil")
	}
	if err.Error() != "connector not found" {
		t.Fatalf("Expected error 'connector not found', got: %v", err)
	}
}

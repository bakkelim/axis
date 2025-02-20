package controllers

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteContract_Success(t *testing.T) {
	const testID = "test_contract"
	filename := filepath.Join(contractsDir, testID+".json")

	// Ensure the contracts directory exists.
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		t.Fatalf("error creating contracts dir: %v", err)
	}

	// Create a dummy contract file.
	if err := os.WriteFile(filename, []byte(`{"dummy": "data"}`), 0644); err != nil {
		t.Fatalf("error writing dummy contract: %v", err)
	}

	// Call deleteContract.
	if err := deleteContract(testID); err != nil {
		t.Fatalf("deleteContract returned error: %v", err)
	}

	// Ensure the file was removed.
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("expected file to be deleted, but still exists: %v", err)
	}
}

func TestStorageDeleteContract_NotFound(t *testing.T) {
	const nonExistentID = "non_existent_contract"
	err := deleteContract(nonExistentID)
	if err == nil {
		t.Fatal("expected an error for non-existent contract, got nil")
	}
	if !errors.Is(err, errors.New("contract not found")) && err.Error() != "contract not found" {
		t.Fatalf("expected error 'contract not found', got: %v", err)
	}
}

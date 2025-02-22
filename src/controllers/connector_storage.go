package controllers

import (
	"axis/src/models"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultConnectorsDir = "../connectors"

var testConnectorsDir string

func getConnectorsDir() string {
	if testConnectorsDir != "" {
		return testConnectorsDir
	}
	return defaultConnectorsDir
}

var saveConnector = func(connector *models.Connector) error {
	// Ensure directory exists before saving
	if err := os.MkdirAll(getConnectorsDir(), 0755); err != nil {
		return err
	}

	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}

	filename := filepath.Join(getConnectorsDir(), connector.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

// LoadConnector retrieves a connector by its ID
func LoadConnector(id string) (*models.Connector, error) {
	filename := filepath.Join(getConnectorsDir(), id+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("connector not found")
		}
		return nil, err
	}

	var connector models.Connector
	if err := json.Unmarshal(data, &connector); err != nil {
		return nil, err
	}
	return &connector, nil
}

func listConnectors() ([]models.Connector, error) {
	files, err := os.ReadDir(getConnectorsDir())
	if err != nil {
		return []models.Connector{}, err
	}

	var connectors = []models.Connector{}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5] // remove .json
		connector, err := LoadConnector(id)
		if err != nil {
			continue
		}
		connectors = append(connectors, *connector)
	}
	return connectors, nil
}

func deleteConnector(id string) error {
	filename := filepath.Join(getConnectorsDir(), id+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return errors.New("connector not found")
		}
		return err
	}
	return nil
}

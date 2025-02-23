package controllers

import (
	"axis/src/models"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const connectorsDir = "../connectors"

func init() {
	// Ensure connectors directory exists
	if err := os.MkdirAll(connectorsDir, 0755); err != nil {
		panic(err)
	}
}

var saveConnector = func(connector *models.Connector) error {
	data, err := json.Marshal(connector)
	if err != nil {
		return err
	}

	filename := filepath.Join(connectorsDir, connector.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

var loadConnector = func (id string) (*models.Connector, error) {
	filename := filepath.Join(connectorsDir, id+".json")
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
	files, err := os.ReadDir(connectorsDir)
	if err != nil {
		return []models.Connector{}, err
	}

	var connectors = []models.Connector{}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5] // remove .json
		connector, err := loadConnector(id)
		if err != nil {
			continue
		}
		connectors = append(connectors, *connector)
	}
	return connectors, nil
}

func deleteConnector(id string) error {
	filename := filepath.Join(connectorsDir, id+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return errors.New("connector not found")
		}
		return err
	}
	return nil
}

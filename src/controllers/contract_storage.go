package controllers

import (
	"axis/src/models"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const defaultContractsDir = "../data-contracts"

var testContractsDir string

func getContractsDir() string {
	if testContractsDir != "" {
		return testContractsDir
	}
	return defaultContractsDir
}

func init() {
	// Ensure contracts directory exists
	if err := os.MkdirAll(getContractsDir(), 0755); err != nil {
		panic(err)
	}
}

var saveContract = func(contract *models.Contract) error {
	data, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	filename := filepath.Join(getContractsDir(), contract.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

func loadContract(id string) (*models.Contract, error) {
	filename := filepath.Join(getContractsDir(), id+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("contract not found")
		}
		return nil, err
	}

	var contract models.Contract
	if err := json.Unmarshal(data, &contract); err != nil {
		return nil, err
	}
	return &contract, nil
}

var listContracts = func() ([]models.Contract, error) {
	files, err := os.ReadDir(getContractsDir())
	if err != nil {
		return []models.Contract{}, err
	}

	var contracts = []models.Contract{}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		id := file.Name()[:len(file.Name())-5] // remove .json
		contract, err := loadContract(id)
		if err != nil {
			continue
		}
		contracts = append(contracts, *contract)
	}
	return contracts, nil
}

func deleteContract(id string) error {
	filename := filepath.Join(getContractsDir(), id+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return errors.New("contract not found")
		}
		return err
	}
	return nil
}

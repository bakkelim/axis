package controllers

import (
	"axis/src/models"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const contractsDir = "../data-contracts"

func init() {
	// Ensure contracts directory exists
	if err := os.MkdirAll(contractsDir, 0755); err != nil {
		panic(err)
	}
}

var saveContract = func(contract *models.Contract) error {
	data, err := json.Marshal(contract)
	if err != nil {
		return err
	}

	filename := filepath.Join(contractsDir, contract.ID+".json")
	return os.WriteFile(filename, data, 0644)
}

func loadContract(id string) (*models.Contract, error) {
	filename := filepath.Join(contractsDir, id+".json")
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

func listContracts() ([]models.Contract, error) {
	files, err := os.ReadDir(contractsDir)
	if err != nil {
		return nil, err
	}

	var contracts []models.Contract
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
	filename := filepath.Join(contractsDir, id+".json")
	if err := os.Remove(filename); err != nil {
		if os.IsNotExist(err) {
			return errors.New("contract not found")
		}
		return err
	}
	return nil
}

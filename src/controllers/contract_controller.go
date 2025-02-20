package controllers

import (
	"axis/src/models"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bytes"
	"database/sql"
	"html/template"
	"time"

	_ "github.com/lib/pq"
)

// CreateContract creates a new contract
func CreateContract(c *gin.Context) {
	var contract models.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract.ID = uuid.New().String()

	if err := saveContract(&contract); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save contract"})
		return
	}

	c.JSON(http.StatusCreated, contract)
}

// ListContracts returns all contracts
func ListContracts(c *gin.Context) {
	contracts, err := listContracts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list contracts"})
		return
	}

	c.JSON(http.StatusOK, contracts)
}

// GetContractByID retrieves a contract by its ID
func GetContractByID(c *gin.Context) {
	id := c.Param("id")

	contract, err := loadContract(id)
	if err != nil {
		if err.Error() == "contract not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contract"})
		}
		return
	}

	c.JSON(http.StatusOK, contract)
}

// UpdateContract updates an existing contract
func UpdateContract(c *gin.Context) {
	id := c.Param("id")

	// Check if contract exists
	if _, err := loadContract(id); err != nil {
		if err.Error() == "contract not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contract"})
		}
		return
	}

	var contract models.Contract
	if err := c.ShouldBindJSON(&contract); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contract.ID = id
	if err := saveContract(&contract); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contract"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

// DeleteContract removes a contract
func DeleteContract(c *gin.Context) {
	id := c.Param("id")

	if err := deleteContract(id); err != nil {
		if err.Error() == "contract not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contract"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contract deleted successfully"})
}

func ExecuteContract(c *gin.Context) {
	id := c.Param("id")

	// 1. Read the contract file
	contractPath := filepath.Join(contractsDir, id+".json")
	contractData, err := os.ReadFile(contractPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	var contract models.Contract
	if err := json.Unmarshal(contractData, &contract); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing contract data"})
		return
	}

	// 2. Load the connector
	connectorPath := filepath.Join("../connectors", fmt.Sprintf("%s.json", contract.Query.ConnectorID))
	connectorData, err := os.ReadFile(connectorPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connector not found"})
		return
	}

	var connector models.Connector
	if err := json.Unmarshal(connectorData, &connector); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing connector data"})
		return
	}

	// 3. Execute the SQL query
	db, err := sql.Open(connector.Type, connector.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	rows, err := db.Query(contract.Query.SQLQuery)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution failed"})
		return
	}
	defer rows.Close()

	// 3. Generate the response template
	var results []map[string]interface{}
	columns, _ := rows.Columns()

	for rows.Next() {
		// Create a slice of interface{} to store the values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			row[col] = val
		}

		results = append(results, row)
	}

	// parse result into template
	parsedResults := make([]map[string]interface{}, 0)
	for _, result := range results {
		// Convert template to string representation
		templateStr := make(map[string]interface{})
		for key, value := range contract.ResponseTemplate.Template {
			if tmpl, ok := value.(string); ok {
				// Create a new template
				t, err := template.New("field").Parse(tmpl)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Template parsing failed"})
					return
				}

				// Execute template with result data
				var buf bytes.Buffer
				if err := t.Execute(&buf, result); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Template execution failed"})
					return
				}
				templateStr[key] = buf.String()
			}
		}
		parsedResults = append(parsedResults, templateStr)
	}

	// 5. Return the response
	response := gin.H{
		"contract_id": id,
		"status":      "success",
		"results":     parsedResults,
		"timestamp":   time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

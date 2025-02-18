package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"axis/src/models"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

// GetContractByID retrieves a contract by its ID
func GetContractByID(c *gin.Context) {
	id := c.Param("id")

	// Read the contract file
	filePath := filepath.Join("../data-contracts", fmt.Sprintf("contract-%s.json", id))
	contractData, err := os.ReadFile(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		return
	}

	var contract models.Contract
	if err := json.Unmarshal(contractData, &contract); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing contract data"})
		return
	}

	c.JSON(http.StatusOK, contract)
}

func ExecuteContract(c *gin.Context) {
	id := c.Param("id")

	// 1. Read the contract file
	contractPath := filepath.Join("../data-contracts", fmt.Sprintf("contract-%s.json", id))
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

	// 2. Execute the SQL query
	connector := contract.Connector
	db, err := sql.Open(connector.Type, connector.ConnectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	rows, err := db.Query(connector.SQLQuery)
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

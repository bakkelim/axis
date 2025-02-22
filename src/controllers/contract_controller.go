package controllers

import (
	"axis/src/models"
	"crypto/sha256"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"bytes"
	"database/sql"
	"html/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// For testing purposes
var sqlOpen = sql.Open

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

	// Parse request body
	var req models.ExecuteContractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Load contract
	contract, err := loadContract(id)
	if err != nil {
		if err.Error() == "contract not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Contract not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load contract"})
		}
		return
	}

	// Apply filters, pagination, and sorting from request if provided
	if req.Filters != nil {
		contract.Query.Filters = req.Filters
	}
	if req.Pagination != nil {
		contract.Query.Pagination = req.Pagination
	}
	if req.Sort != nil {
		contract.Query.Sort = req.Sort
	}

	// Load the connector
	connector, err := LoadConnector(contract.Query.ConnectorID)
	if err != nil {
		if err.Error() == "connector not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connector not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load connector"})
		}
		return
	}

	// Execute the SQL query
	db, err := sqlOpen(connector.Type, buildConnectionString(connector.Config, connector.Type))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	baseQuery := contract.Query.SQLQuery
	whereClause, values := buildWhereClause(contract.Query.Filters)
	orderByClause := buildOrderByClause(contract.Query.Sort)
	query := baseQuery + whereClause + orderByClause

	if contract.Query.Pagination != nil {
		offset := (contract.Query.Pagination.Page - 1) * contract.Query.Pagination.PageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d",
			contract.Query.Pagination.PageSize, offset)
	}

	rows, err := db.Query(query, values...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Query execution failed"})
		return
	}
	defer rows.Close()

	// Generate the response template
	var results []map[string]interface{}
	columns, _ := rows.Columns()

	for rows.Next() {
		// Create a map for this row's data
		rowData := make(map[string]interface{})

		// Create properly typed containers for the scan
		scanArgs := make([]interface{}, len(columns))
		for i := range columns {
			scanArgs[i] = new(interface{})
		}

		if err := rows.Scan(scanArgs...); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning row"})
			return
		}

		// Copy the results into the row map
		for i, col := range columns {
			val := *(scanArgs[i].(*interface{}))
			// Convert []byte to string for MySQL text-based columns
			if b, ok := val.([]byte); ok {
				rowData[col] = string(b)
			} else {
				rowData[col] = val
			}
		}

		results = append(results, rowData)
	}

	// parse result into template
	parsedResults := make([]map[string]interface{}, 0)
	for _, result := range results {
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

				// Check if this field needs anonymization
				fieldValue := buf.String()
				for _, rule := range contract.ResponseTemplate.Anonymization {
					if rule.Field == key {
						fieldValue = anonymizeValue(fieldValue, rule)
						break
					}
				}
				templateStr[key] = fieldValue
			}
		}
		parsedResults = append(parsedResults, templateStr)
	}

	// Return the response
	response := gin.H{
		"contract_id": id,
		"status":      "success",
		"results":     parsedResults,
		"timestamp":   time.Now().UTC(),
	}

	c.JSON(http.StatusOK, response)
}

// Add these helper functions before ExecuteContract
func anonymizeValue(value string, rule models.AnonymizationRule) string {
	switch rule.Method {
	case "mask":
		if rule.Pattern != "" {
			return applyMaskPattern(value, rule.Pattern)
		}
		// Default mask if no pattern provided
		return strings.Repeat("*", len(value))
	case "hash":
		h := sha256.New()
		h.Write([]byte(value))
		return fmt.Sprintf("%x", h.Sum(nil))
	case "randomize":
		return uuid.New().String() // Simple randomization
	default:
		return value
	}
}

func applyMaskPattern(value string, pattern string) string {
	result := []rune(pattern)
	valueRunes := []rune(value)
	valueIndex := 0

	for i, r := range result {
		if r == 'X' {
			if valueIndex < len(valueRunes) {
				// Skip any dashes in the input value
				for valueIndex < len(valueRunes) && valueRunes[valueIndex] == '-' {
					valueIndex++
				}
				if valueIndex < len(valueRunes) {
					result[i] = valueRunes[valueIndex]
					valueIndex++
				} else {
					result[i] = '*'
				}
			} else {
				result[i] = '*'
			}
		}
	}
	return string(result)
}

// Parse the database configuration and build the connection string
func buildConnectionString(config models.DatabaseConfig, connectorType string) string {
	switch connectorType {
	case "postgres":
		return fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.DBName,
		)
	case "mysql":
		connString := fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s",
			config.User, config.Password, config.Host, config.Port, config.DBName,
		)
		return connString
	default:
		return ""
	}
}

func buildWhereClause(filters []models.FilterCondition) (string, []interface{}) {
	if len(filters) == 0 {
		return "", nil
	}

	var conditions []string
	var values []interface{}
	paramCount := 1

	for _, filter := range filters {
		switch filter.Operator {
		case models.OperatorEquals:
			conditions = append(conditions, fmt.Sprintf("%s = $%d", filter.Field, paramCount))
			values = append(values, filter.Value)
			paramCount++
		case models.OperatorNotEquals:
			conditions = append(conditions, fmt.Sprintf("%s != $%d", filter.Field, paramCount))
			values = append(values, filter.Value)
			paramCount++
		case models.OperatorGreater:
			conditions = append(conditions, fmt.Sprintf("%s > $%d", filter.Field, paramCount))
			values = append(values, filter.Value)
			paramCount++
		case models.OperatorLess:
			conditions = append(conditions, fmt.Sprintf("%s < $%d", filter.Field, paramCount))
			values = append(values, filter.Value)
			paramCount++
		case models.OperatorLike:
			conditions = append(conditions, fmt.Sprintf("%s LIKE $%d", filter.Field, paramCount))
			values = append(values, filter.Value)
			paramCount++
		case models.OperatorIn:
			if inValues, ok := filter.Value.([]interface{}); ok {
				placeholders := make([]string, len(inValues))
				for i := range inValues {
					placeholders[i] = fmt.Sprintf("$%d", paramCount+i)
				}
				conditions = append(conditions, fmt.Sprintf("%s IN (%s)",
					filter.Field, strings.Join(placeholders, ",")))
				values = append(values, inValues...)
				paramCount += len(inValues)
			} else {
				// If not a valid slice, skip this condition but keep the value
				values = append(values, filter.Value)
			}
		}
	}

	return " WHERE " + strings.Join(conditions, " AND "), values
}

func buildOrderByClause(sortOptions []models.SortOption) string {
	if len(sortOptions) == 0 {
		return ""
	}

	var orderByClauses []string
	for _, sortOption := range sortOptions {
		orderByClauses = append(orderByClauses, fmt.Sprintf("%s %s", sortOption.Field, sortOption.Direction))
	}

	return " ORDER BY " + strings.Join(orderByClauses, ", ")
}

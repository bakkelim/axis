package models

// DatabaseConfig represents the database connection configuration
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DBName   string `json:"dbname"`
}

// Connector represents the database connection configuration
type Connector struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Config      DatabaseConfig `json:"config"`
}

// FilterOperator represents the type of filter operation
type FilterOperator string

const (
	OperatorEquals    FilterOperator = "eq"
	OperatorNotEquals FilterOperator = "neq"
	OperatorGreater   FilterOperator = "gt"
	OperatorLess      FilterOperator = "lt"
	OperatorLike      FilterOperator = "like"
	OperatorIn        FilterOperator = "in"
)

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    string         `json:"field"`    // Column name to filter on
	Operator FilterOperator `json:"operator"` // Filter operation to apply
	Value    any            `json:"value"`    // Value to compare against
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page     int `json:"page"`     // Page number (1-based)
	PageSize int `json:"pageSize"` // Number of items per page
}

// SortOption represents a single sorting option
type SortOption struct {
	Field     string `json:"field"`     // Column name to sort by
	Direction string `json:"direction"` // Sort direction ("asc" or "desc")
}

// DatabaseQuery represents the query configuration
type DatabaseQuery struct {
	ConnectorID string             `json:"connectorId"`
	SQLQuery    string             `json:"sqlQuery"`
	Filters     []FilterCondition  `json:"filters,omitempty"`
	Pagination  *PaginationOptions `json:"pagination,omitempty"`
	Sort        []SortOption       `json:"sort,omitempty"`
}

// ExecuteContractRequest represents the request body for contract execution
type ExecuteContractRequest struct {
	Filters    []FilterCondition  `json:"filters,omitempty"`
	Pagination *PaginationOptions `json:"pagination,omitempty"`
	Sort       []SortOption       `json:"sort,omitempty"`
}

// AnonymizationRule defines how a field should be anonymized
type AnonymizationRule struct {
	Field   string `json:"field"`   // The field name to anonymize
	Method  string `json:"method"`  // The anonymization method (e.g., "mask", "hash", "randomize")
	Pattern string `json:"pattern"` // Optional pattern for masking (e.g., "XXX-XX-****" for SSN)
}

// ResponseTemplate represents the template structure for API responses
type ResponseTemplate struct {
	ID            string              `json:"id"`
	Template      map[string]any      `json:"template"`
	Anonymization []AnonymizationRule `json:"anonymization,omitempty"`
}

// Contract represents the main contract structure
type Contract struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Query            DatabaseQuery    `json:"query"`
	ResponseTemplate ResponseTemplate `json:"responseTemplate"`
}

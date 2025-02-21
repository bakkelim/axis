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

// DatabaseQuery represents the query configuration
type DatabaseQuery struct {
	ConnectorID string `json:"connectorId"`
	SQLQuery    string `json:"sqlQuery"`
}

// AnonymizationRule defines how a field should be anonymized
type AnonymizationRule struct {
	Field   string `json:"field"`   // The field name to anonymize
	Method  string `json:"method"`  // The anonymization method (e.g., "mask", "hash", "randomize")
	Pattern string `json:"pattern"` // Optional pattern for masking (e.g., "XXX-XX-****" for SSN)
}

// ResponseTemplate represents the template structure for API responses
type ResponseTemplate struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Template      map[string]interface{} `json:"template"`
	Anonymization []AnonymizationRule    `json:"anonymization,omitempty"`
}

// Contract represents the main contract structure
type Contract struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Query            DatabaseQuery    `json:"query"`
	ResponseTemplate ResponseTemplate `json:"responseTemplate"`
}

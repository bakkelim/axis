package models

// Connector represents the database connection configuration
type Connector struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	ConnectionString string `json:"connectionString"`
}

// DatabaseQuery represents the query configuration
type DatabaseQuery struct {
	ConnectorID string `json:"connectorId"`
	SQLQuery    string `json:"sqlQuery"`
}

// ResponseTemplate represents the template structure for API responses
type ResponseTemplate struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Template    map[string]interface{} `json:"template"`
}

// Contract represents the main contract structure
type Contract struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	Description      string           `json:"description"`
	Query            DatabaseQuery    `json:"query"`
	ResponseTemplate ResponseTemplate `json:"responseTemplate"`
}

package models

// Connector represents the database connection configuration
type Connector struct {
    ID               int    `json:"id"`
    Name             string `json:"name"`
    Description      string `json:"description"`
    Type             string `json:"type"`
    ConnectionString string `json:"connectionString"`
    SQLQuery         string `json:"SQLQuery"`
}

// ResponseTemplate represents the template structure for API responses
type ResponseTemplate struct {
    ID          int                    `json:"id"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Template    map[string]interface{} `json:"template"`
}

// Contract represents the main contract structure
type Contract struct {
    ID               int              `json:"id"`
    Name             string           `json:"name"`
    Description      string           `json:"description"`
    Connector        Connector        `json:"connector"`
    ResponseTemplate ResponseTemplate `json:"responseTemplate"`
}
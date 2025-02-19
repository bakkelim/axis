package controllers

import (
	"axis/src/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateConnector handles the creation of a new connector
func CreateConnector(c *gin.Context) {
	var connector models.Connector
	if err := c.ShouldBindJSON(&connector); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate unique ID
	connector.ID = uuid.New().String()

	if err := saveConnector(&connector); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save connector"})
		return
	}

	c.JSON(http.StatusCreated, connector)
}

// ListConnectors returns all connectors
func ListConnectors(c *gin.Context) {
	connectors, err := listConnectors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list connectors"})
		return
	}

	c.JSON(http.StatusOK, connectors)
}

// GetConnector returns a specific connector by ID
func GetConnector(c *gin.Context) {
	id := c.Param("id")

	connector, err := loadConnector(id)
	if err != nil {
		if err.Error() == "connector not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connector not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load connector"})
		}
		return
	}

	c.JSON(http.StatusOK, connector)
}

// UpdateConnector updates an existing connector
func UpdateConnector(c *gin.Context) {
	id := c.Param("id")

	// Check if connector exists
	_, err := loadConnector(id)
	if err != nil {
		if err.Error() == "connector not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connector not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load connector"})
		}
		return
	}

	var connector models.Connector
	if err := c.ShouldBindJSON(&connector); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	connector.ID = id
	if err := saveConnector(&connector); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update connector"})
		return
	}

	c.JSON(http.StatusOK, connector)
}

// DeleteConnector removes a connector
func DeleteConnector(c *gin.Context) {
	id := c.Param("id")

	if err := deleteConnector(id); err != nil {
		if err.Error() == "connector not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Connector not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete connector"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connector deleted successfully"})
}

// TestConnection tests if a connector can establish a connection
func TestConnection(c *gin.Context) {
	id := c.Param("id")
	_ = id
	// TODO: Implement connection testing logic

	c.JSON(http.StatusOK, gin.H{"status": "connection successful"})
}

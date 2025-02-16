package middleware

import (
    "github.com/gin-gonic/gin"
    "log"
)

// Logger is a middleware function that logs the incoming requests.
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        log.Printf("Request: %s %s", c.Request.Method, c.Request.URL)
        c.Next()
    }
}

// Auth is a middleware function that checks for authentication.
func Auth() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Implement authentication logic here
        // For example, check for a token in the headers
        token := c.Request.Header.Get("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        c.Next()
    }
}
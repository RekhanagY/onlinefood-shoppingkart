package main

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func AuthApiKey() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Authentication is optional as per requirement.
        apiKey := c.GetHeader("api_key")
        if apiKey == "" {
            c.Next()
            return
        }
        if apiKey != "apitest" {
            c.AbortWithStatusJSON(401, gin.H{
                "code":    http.StatusUnauthorized,
                "message": http.StatusText(http.StatusUnauthorized),
            })
            return
        }
        c.Next()
    }
}

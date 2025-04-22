package middleware

import (
	"github.com/gin-gonic/gin"
	"time"
)

// Logger é um middleware que registra o tempo de resposta das requisições
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Processa a requisição
		c.Next()

		// Calcula o tempo de resposta
		duration := time.Since(start)

		// Log da requisição
		c.JSON(c.Writer.Status(), gin.H{
			"status":   c.Writer.Status(),
			"method":   c.Request.Method,
			"path":     c.Request.URL.Path,
			"duration": duration.String(),
		})
	}
}

// CORS configura os headers de CORS
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

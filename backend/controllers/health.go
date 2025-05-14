package controllers

import (
	"competition-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthCheck handles requests to the health check endpoint
func HealthCheck(c *gin.Context) {
	// Check database connection
	dbStatus := "up"
	if err := models.DB.Ping(); err != nil {
		dbStatus = "down"
	}

	// Check Redis connection
	redisStatus := "up"
	_, err := models.RedisClient.Ping(models.GetContext()).Result()
	if err != nil {
		redisStatus = "down"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"database": dbStatus,
		"redis":    redisStatus,
	})
}

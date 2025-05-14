package controllers

import (
	"competition-app/models"
	"competition-app/validation"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetCompetitions handles requests to get all competitions
func GetCompetitions(c *gin.Context) {
	// Try to get data from cache first
	cacheKey := "competitions:all"
	cachedData, err := models.GetCache(cacheKey)
	
	if err == nil {
		var competitions []models.Competition
		if err := json.Unmarshal([]byte(cachedData), &competitions); err == nil {
			c.JSON(http.StatusOK, competitions)
			return
		}
	}

	// If cache miss, get from database
	competitions, err := models.GetAllCompetitions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve competitions", "details": err.Error()})
		return
	}

	// Cache the result
	if competitionsJson, err := json.Marshal(competitions); err == nil {
		models.SetCache(cacheKey, string(competitionsJson), 5*time.Minute)
	}

	c.JSON(http.StatusOK, competitions)
}

// GetCompetition handles requests to get a specific competition
func GetCompetition(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
		return
	}

	// Try to get data from cache first
	cacheKey := "competitions:" + strconv.Itoa(id)
	cachedData, err := models.GetCache(cacheKey)
	
	if err == nil && cachedData != "" {
		var competition models.Competition
		if err := json.Unmarshal([]byte(cachedData), &competition); err == nil {
			c.JSON(http.StatusOK, competition)
			return
		}
	}

	// If cache miss, get from database
	competition, err := models.GetCompetition(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Cache the result
	if competitionJson, err := json.Marshal(competition); err == nil {
		models.SetCache(cacheKey, string(competitionJson), 5*time.Minute)
	}

	c.JSON(http.StatusOK, competition)
}

// CreateCompetition handles requests to create a new competition
func CreateCompetition(c *gin.Context) {
	var competition models.Competition

	// Bind the request body to the competition struct
	if err := c.ShouldBindJSON(&competition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Convert to validation type for validation
	validationObj := validation.Competition{
		Name:        competition.Name,
		Description: competition.Description,
		Date:        competition.Date,
		Location:    competition.Location,
	}

	// Validate the competition data
	if err := validation.ValidateCompetition(&validationObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the competition
	if err := models.CreateCompetition(&competition); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create competition", "details": err.Error()})
		return
	}

	// Invalidate cache
	models.DeleteCache("competitions:all")

	// Cache the newly created competition
	if competitionJson, err := json.Marshal(competition); err == nil {
		models.SetCache("competitions:"+strconv.Itoa(competition.ID), string(competitionJson), 5*time.Minute)
	}

	c.JSON(http.StatusCreated, competition)
}

// UpdateCompetition handles requests to update an existing competition
func UpdateCompetition(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
		return
	}

	var competition models.Competition

	// Bind the request body to the competition struct
	if err := c.ShouldBindJSON(&competition); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Set the ID from the URL parameter
	competition.ID = id

	// Convert to validation type for validation
	validationObj := validation.Competition{
		ID:          competition.ID,
		Name:        competition.Name,
		Description: competition.Description,
		Date:        competition.Date,
		Location:    competition.Location,
	}

	// Validate the competition data
	if err := validation.ValidateCompetition(&validationObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the competition
	if err := models.UpdateCompetition(&competition); err != nil {
		if err.Error() == "competition not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update competition", "details": err.Error()})
		}
		return
	}

	// Invalidate cache
	models.DeleteCache("competitions:all")
	models.DeleteCache("competitions:" + strconv.Itoa(id))

	c.JSON(http.StatusOK, competition)
}

// DeleteCompetition handles requests to delete a competition
func DeleteCompetition(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
		return
	}

	// Delete the competition
	if err := models.DeleteCompetition(id); err != nil {
		if err.Error() == "competition not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete competition", "details": err.Error()})
		}
		return
	}

	// Invalidate cache
	models.DeleteCache("competitions:all")
	models.DeleteCache("competitions:" + strconv.Itoa(id))
	models.DeleteCache("participants:competition:" + strconv.Itoa(id))

	c.JSON(http.StatusOK, gin.H{"message": "Competition deleted successfully"})
}

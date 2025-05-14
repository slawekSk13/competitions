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

// GetParticipants handles requests to get all participants
func GetParticipants(c *gin.Context) {
	// Check if we're filtering by competition ID
	competitionIDStr := c.Query("competition_id")

	// Try to get data from cache first
	var cacheKey string
	var err error
	if competitionIDStr != "" {
		cacheKey = "participants:competition:" + competitionIDStr
	} else {
		cacheKey = "participants:all"
	}

	var cachedData string
	cachedData, err = models.GetCache(cacheKey)
	if err == nil {
		var participants []models.Participant
		err = json.Unmarshal([]byte(cachedData), &participants)
		if err == nil {
			c.JSON(http.StatusOK, participants)
			return
		}
	}

	// If cache miss, get from database
	var participants []models.Participant

	if competitionIDStr != "" {
		var competitionID int
		competitionID, err = strconv.Atoi(competitionIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
			return
		}

		participants, err = models.GetParticipantsByCompetition(competitionID)
	} else {
		participants, err = models.GetAllParticipants()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve participants", "details": err.Error()})
		return
	}

	// Cache the result
	var participantsJson []byte
	participantsJson, err = json.Marshal(participants)
	if err == nil {
		models.SetCache(cacheKey, string(participantsJson), 5*time.Minute)
	}

	c.JSON(http.StatusOK, participants)
}

// GetParticipant handles requests to get a specific participant
func GetParticipant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
		return
	}

	// Try to get data from cache first
	cacheKey := "participants:" + strconv.Itoa(id)
	cachedData, err := models.GetCache(cacheKey)
	if err == nil {
		var participant models.Participant
		err = json.Unmarshal([]byte(cachedData), &participant)
		if err == nil {
			c.JSON(http.StatusOK, participant)
			return
		}
	}

	participant, err := models.GetParticipant(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Cache the result
	if participantJson, err := json.Marshal(participant); err == nil {
		models.SetCache(cacheKey, string(participantJson), 5*time.Minute)
	}

	c.JSON(http.StatusOK, participant)
}

// GetParticipantCompetitions handles requests to get all competitions for a participant
func GetParticipantCompetitions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
		return
	}

	competitions, err := models.GetParticipantCompetitions(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve competitions", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, competitions)
}

// CreateParticipant handles requests to create a new participant
func CreateParticipant(c *gin.Context) {
	var participant models.Participant

	if err := c.ShouldBindJSON(&participant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	// Convert to validation type for validation
	validationObj := validation.Participant{
		Name:  participant.Name,
		Email: participant.Email,
	}

	if err := validation.ValidateParticipant(&validationObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.CreateParticipant(&participant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create participant", "details": err.Error()})
		return
	}

	// Invalidate cache
	models.DeleteCache("participants:all")

	c.JSON(http.StatusCreated, participant)
}

// AddParticipantToCompetition handles requests to add a participant to a competition
func AddParticipantToCompetition(c *gin.Context) {
	participantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
		return
	}

	var data struct {
		CompetitionID    int    `json:"competition_id" binding:"required"`
		RegistrationDate string `json:"registration_date" binding:"required"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}

	registrationDate, err := time.Parse("2006-01-02", data.RegistrationDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid registration date format"})
		return
	}

	if err := models.AddParticipantToCompetition(participantID, data.CompetitionID, registrationDate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add participant to competition", "details": err.Error()})
		return
	}

	// Invalidate cache
	models.DeleteCache("participants:all")
	models.DeleteCache("participants:competition:" + strconv.Itoa(data.CompetitionID))

	c.JSON(http.StatusOK, gin.H{"message": "Participant added to competition successfully"})
}

// UpdateParticipant handles requests to update an existing participant
func UpdateParticipant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
		return
	}

	var participant models.Participant
	if err := c.ShouldBindJSON(&participant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format", "details": err.Error()})
		return
	}
	participant.ID = id

	validationObj := validation.Participant{
		ID:    participant.ID,
		Name:  participant.Name,
		Email: participant.Email,
	}

	if err := validation.ValidateParticipant(&validationObj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.UpdateParticipant(&participant); err != nil {
		if err.Error() == "participant not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update participant", "details": err.Error()})
		}
		return
	}

	// Invalidate cache
	models.DeleteCache("participants:all")
	models.DeleteCache("participants:" + strconv.Itoa(id))

	c.JSON(http.StatusOK, participant)
}

// RemoveParticipantFromCompetition handles requests to remove a participant from a competition
func RemoveParticipantFromCompetition(c *gin.Context) {
	participantID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
		return
	}

	competitionID, err := strconv.Atoi(c.Param("competition_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid competition ID"})
		return
	}

	if err := models.RemoveParticipantFromCompetition(participantID, competitionID); err != nil {
		if err.Error() == "participant not found in competition" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove participant from competition", "details": err.Error()})
		}
		return
	}

	// Invalidate cache
	models.DeleteCache("participants:all")
	models.DeleteCache("participants:competition:" + strconv.Itoa(competitionID))

	c.JSON(http.StatusOK, gin.H{"message": "Participant removed from competition successfully"})
}

// DeleteParticipant handles requests to delete a participant
func DeleteParticipant(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant ID"})
		return
	}

	if err := models.DeleteParticipant(id); err != nil {
		if err.Error() == "participant not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete participant", "details": err.Error()})
		}
		return
	}

	// Invalidate cache
	models.DeleteCache("participants:all")
	models.DeleteCache("participants:" + strconv.Itoa(id))

	c.JSON(http.StatusOK, gin.H{"message": "Participant deleted successfully"})
}

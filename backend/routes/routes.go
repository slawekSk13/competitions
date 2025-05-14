package routes

import (
	"competition-app/controllers"
	"competition-app/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	// Healthcheck
	router.GET("/api", controllers.HealthCheck)

	// Competitions API
	competitions := router.Group("/api/competitions")
	{
		competitions.GET("", controllers.GetCompetitions)
		competitions.GET("/:id", controllers.GetCompetition)
		competitions.POST("", controllers.CreateCompetition)
		competitions.PUT("/:id", controllers.UpdateCompetition)
		competitions.DELETE("/:id", controllers.DeleteCompetition)
	}

	// Participants API
	participants := router.Group("/api/participants")
	{
		participants.GET("", controllers.GetParticipants)
		participants.GET("/:id", controllers.GetParticipant)
		participants.GET("/:id/competitions", controllers.GetParticipantCompetitions)
		participants.POST("", controllers.CreateParticipant)
		participants.POST("/:id/competitions", controllers.AddParticipantToCompetition)
		participants.PUT("/:id", controllers.UpdateParticipant)
		participants.DELETE("/:id/competitions/:competition_id", controllers.RemoveParticipantFromCompetition)
		participants.DELETE("/:id", controllers.DeleteParticipant)
	}

	return router
}

package routes

import (
	"github.com/codepnw/ticket-api/handlers"
	eventRepository "github.com/codepnw/ticket-api/repositories/event"
	eventService "github.com/codepnw/ticket-api/services/event"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func EventRoutes(db *sqlx.DB, r *gin.Engine, version string) {
	eventRepo := eventRepository.NewEventRepository(db)
	eventSrv := eventService.NewEventService(eventRepo)
	eventHandler := handlers.NewEventHandler(eventSrv)

	g := r.Group(version + "/events")

	g.GET("/", eventHandler.GetEvents)
	g.GET("/:event_id", eventHandler.GetEvent)
	g.POST("/", eventHandler.NewEvent)
	g.PATCH("/:event_id", eventHandler.UpdateEvent)
	g.DELETE("/:event_id", eventHandler.DeleteEvent)
}

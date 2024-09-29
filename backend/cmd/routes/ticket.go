package routes

import (
	"github.com/codepnw/ticket-api/handlers"
	ticketRepository "github.com/codepnw/ticket-api/repositories/ticket"
	ticketService "github.com/codepnw/ticket-api/services/ticket"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func TicketRoutes(db *sqlx.DB, r *gin.Engine, version string) {
	ticketRepo := ticketRepository.NewTicketRepository(db)
	ticketSrv := ticketService.NewTicketService(ticketRepo)
	ticketHandler := handlers.NewTicketHandler(ticketSrv)

	g := r.Group(version + "/tickets")

	g.GET("/", ticketHandler.GetTickets)
	g.GET("/:ticket_id", ticketHandler.GetTicket)
	g.POST("/", ticketHandler.NewTicket)
	g.PATCH("/:ticket_id", ticketHandler.UpdateTicket)
}

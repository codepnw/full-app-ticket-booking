package handlers

import (
	"net/http"
	"strconv"

	ticketRepository "github.com/codepnw/ticket-api/repositories/ticket"
	ticketService "github.com/codepnw/ticket-api/services/ticket"
	"github.com/gin-gonic/gin"
)

type ticketHandler struct {
	ticketSrv ticketService.ITicketService
}

func NewTicketHandler(ticketSrv ticketService.ITicketService) *ticketHandler {
	return &ticketHandler{ticketSrv: ticketSrv}
}

func (h *ticketHandler) GetTickets(c *gin.Context) {
	tickets, err := h.ticketSrv.GetTickets()
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *ticketHandler) GetTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("ticket_id"))

	ticket, err := h.ticketSrv.GetTicket(uint(id))
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, ticket)
}

func (h *ticketHandler) NewTicket(c *gin.Context) {
	var request ticketService.TicketRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	ticket, err := h.ticketSrv.CreateTicket(request)
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

func (h *ticketHandler) UpdateTicket(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("ticket_id"))
	data := new(ticketRepository.ValidateTicket)

	if err := h.ticketSrv.UpdateTicket(uint(id), data); err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, nil)
}

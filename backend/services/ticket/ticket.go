package ticketService

import (
	eventRepository "github.com/codepnw/ticket-api/repositories/event"
	ticketRepository "github.com/codepnw/ticket-api/repositories/ticket"
)

type ITicketService interface {
	GetTickets() ([]*TicketResponse, error)
	GetTicket(ticketID uint) (*TicketResponse, error)
	CreateTicket(request TicketRequest) (*TicketResponse, error)
	ValidateTicket(ticketID uint, validate *ticketRepository.ValidateTicket) error
}

type TicketResponse struct {
	ID        uint                  `json:"id"`
	EventID   uint                  `json:"eventId"`
	Event     eventRepository.Event `json:"event"`
	Entered   bool                  `json:"entered" default:"false"`
	CreatedAt string                `json:"createdAt"`
	UpdatedAt string                `json:"updatedAt"`
}

type TicketRequest struct {
	EventID uint `json:"eventId"`
	Entered bool `json:"entered" default:"false"`
}

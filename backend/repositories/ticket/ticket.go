package ticketRepository

import (
	"context"

	eventRepository "github.com/codepnw/ticket-api/repositories/event"
)

type ITicketRepository interface {
	GetMany(ctx context.Context) ([]*Ticket, error)
	GetOne(ctx context.Context, ticketID uint) (*Ticket, error)
	CreateOne(ctx context.Context, ticket *Ticket) (*Ticket, error)
	UpdateOne(ctx context.Context, ticketId uint, validate *ValidateTicket) error
}

type Ticket struct {
	ID        uint                  `db:"id" json:"id"`
	EventID   uint                  `db:"event_id" json:"eventId"`
	Event     eventRepository.Event `db:"event" json:"event"`
	Entered   bool                  `db:"entered" default:"false" json:"entered"`
	CreatedAt string                `db:"created_at" json:"createdAt"`
	UpdatedAt string                `db:"updated_at" json:"updatedAt"`
}

type ValidateTicket struct {
	TicketId uint `json:"ticketId"`
	Entered  bool `json:"entered"`
}

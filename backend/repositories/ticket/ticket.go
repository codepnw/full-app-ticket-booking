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
	ID        uint                  `db:"id"`
	EventID   uint                  `db:"event_id"`
	Event     eventRepository.Event `db:"event"`
	Entered   bool                  `db:"entered" default:"false"`
	CreatedAt string                `db:"created_at"`
	UpdatedAt string                `db:"updated_at"`
}

type ValidateTicket struct {
	TicketId uint `json:"ticketId"`
	Entered  bool `json:"entered"`
}

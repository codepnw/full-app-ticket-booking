package ticketRepository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type ticketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) ITicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) GetMany(ctx context.Context) ([]*Ticket, error) {
	query := `
		SELECT id, event_id, event, entered, created_at, updated_at
		FROM tickets;
	`
	tickets := []*Ticket{}

	if err := r.db.Select(&tickets, query); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) GetOne(ctx context.Context, ticketID uint) (*Ticket, error) {
	query := `
		SELECT id, event_id, event, entered, created_at, updated_at
		FROM tickets
		WHERE id = $1;
	`
	ticket := new(Ticket)

	if err := r.db.Get(ticket, query, ticketID); err != nil {
		return nil, err
	}
	return ticket, nil
}

func (r *ticketRepository) CreateOne(ctx context.Context, ticket *Ticket) (*Ticket, error) {
	query := `
		INSERT INTO tickets (event_id, event, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id;
	`
	var lastInsertID int

	err := r.db.QueryRow(query, ticket.EventID, ticket.Event).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	ticket.ID = uint(lastInsertID)
	return ticket, nil
}

func (r *ticketRepository) UpdateOne(ctx context.Context, ticketId uint, validate *ValidateTicket) error {
	query := `
		UPDATE tickets SET entered = $1
		WHERE id = $2;
	`

	_, err := r.db.Exec(query,validate.Entered, ticketId)
	if err != nil {
		return err
	}

	return nil
}

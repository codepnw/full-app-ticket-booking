package eventRepository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type eventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) IEventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) GetMany(ctx context.Context) ([]*Event, error) {
	query := `
		SELECT event_id, name, location, date, created_at, updated_at
		FROM events;
	`
	events := []*Event{}

	if err := r.db.Select(&events, query); err != nil {
		return nil, err
	}
	return events, nil
}

func (r *eventRepository) GetOne(ctx context.Context, eventID string) (*Event, error) {
	query := `
		SELECT event_id, name, location, date, created_at, updated_at
		FROM events
		WHERE event_id = $1;
	`
	event := new(Event)

	if err := r.db.Get(event, query, eventID); err != nil {
		return nil, err
	}
	return event, nil
}

func (r *eventRepository) CreateOne(ctx context.Context, event Event) (*Event, error) {
	query := `
		INSERT INTO events (name, location, date)
		VALUES ($1, $2, $3)
		RETURNING event_id;
	`
	var lastInsertID int

	err := r.db.QueryRow(query, event.Name, event.Location, event.Date).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	event.ID = lastInsertID
	return &event, nil
}

func (r *eventRepository) UpdateOne(ctx context.Context, eventID uint, updateData *EventUpdateRequest) error {
	query := `
		UPDATE events SET name = $1, location = $2, date = $3
		WHERE event_id = $4;
	`
	_, err := r.db.Exec(query, updateData.Name, updateData.Location, updateData.Date, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventRepository) DeleteOne(ctx context.Context, eventID string) error {
	query := `DELETE FROM events WHERE event_id = $1;`

	if _, err := r.db.Exec(query, eventID); err != nil {
		return err
	}
	return nil
}

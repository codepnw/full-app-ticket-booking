package eventRepository

import (
	"context"
)

type IEventRepository interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventID string) (*Event, error)
	CreateOne(ctx context.Context, event Event) (*Event, error)
	UpdateOne(ctx context.Context, eventID uint, updateData *EventUpdateRequest) error
	DeleteOne(ctx context.Context, eventID string) error
}

type Event struct {
	ID        int    `db:"event_id"`
	Name      string `db:"name"`
	Location  string `db:"location"`
	Date      string `db:"date"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type EventUpdateRequest struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	Date      string `json:"date"`
	UpdatedAt string `db:"updated_at"`
}

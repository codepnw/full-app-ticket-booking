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
	ID        int    `db:"event_id" json:"eventId"`
	Name      string `db:"name" json:"name"`
	Location  string `db:"location" json:"location"`
	Date      string `db:"date" json:"date"`
	CreatedAt string `db:"created_at" json:"createdAt"`
	UpdatedAt string `db:"updated_at" json:"updatedAt"`
}

type EventUpdateRequest struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	Date      string `json:"date"`
	UpdatedAt string `json:"updated_at"`
}

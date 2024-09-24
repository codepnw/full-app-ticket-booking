package eventRepository

import (
	"context"
	"time"
)

type IEventRepository interface {
	GetMany(ctx context.Context) ([]*Event, error)
	GetOne(ctx context.Context, eventID string) (*Event, error)
	CreateOne(ctx context.Context, event Event) (*Event, error)
}

type Event struct {
	ID        int       `db:"event_id"`
	Name      string    `db:"name"`
	Location  string    `db:"location"`
	Date      time.Time `db:"date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

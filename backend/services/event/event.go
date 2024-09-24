package eventService

import "time"

type IEventService interface {
	GetEvents() ([]*EventResponse, error)
	GetEvent(eventID string) (*EventResponse, error)
	CreateEvent(request NewEventRequest) (*EventResponse, error)
}

type NewEventRequest struct {
	Name     string    `json:"name"`
	Location string    `json:"location"`
	Date     time.Time `json:"date"`
}

type EventResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Location  string    `json:"location"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

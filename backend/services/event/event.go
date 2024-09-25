package eventService

type IEventService interface {
	GetEvents() ([]*EventResponse, error)
	GetEvent(eventID string) (*EventResponse, error)
	CreateEvent(request EventRequest) (*EventResponse, error)
	UpdateOne(eventID uint, updateData *EventRequest) error
	DeleteOne(eventID string) error
}

type EventRequest struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Date     string `json:"date"`
}

type EventResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	Date      string `json:"date"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

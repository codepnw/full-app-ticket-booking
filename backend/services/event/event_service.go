package eventService

import (
	"context"
	"log"
	"time"

	"github.com/codepnw/ticket-api/pkg/errs"
	eventRepository "github.com/codepnw/ticket-api/repositories/event"
)

type eventService struct {
	eventRepo eventRepository.IEventRepository
}

func NewEventService(eventRepo eventRepository.IEventRepository) IEventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) GetEvents() ([]*EventResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	events, err := s.eventRepo.GetMany(context)
	if err != nil {
		log.Println(err)
		return nil, errs.NewErrUnexpected()
	}

	responses := []*EventResponse{}

	for _, e := range events {
		responses = append(responses, &EventResponse{
			ID:        e.ID,
			Name:      e.Name,
			Location:  e.Location,
			Date:      e.Date,
			CreatedAt: e.CreatedAt,
			UpdatedAt: e.UpdatedAt,
		})
	}
	return responses, nil
}

func (s *eventService) GetEvent(eventID string) (*EventResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	event, err := s.eventRepo.GetOne(context, eventID)
	if err != nil {
		return nil, errs.NewErrUnexpected()
	}

	response := &EventResponse{
		ID:        event.ID,
		Name:      event.Name,
		Location:  event.Location,
		Date:      event.Date,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}

	return response, nil
}

func (s *eventService) CreateEvent(request NewEventRequest) (*EventResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	event := eventRepository.Event{
		Name:     request.Name,
		Location: request.Location,
		Date:     request.Date,
	}

	newEvent, err := s.eventRepo.CreateOne(context, event)
	if err != nil {
		return nil, errs.NewErrUnexpected()
	}

	response := &EventResponse{
		ID:        newEvent.ID,
		Name:      newEvent.Name,
		Location:  newEvent.Location,
		Date:      newEvent.Date,
		CreatedAt: newEvent.CreatedAt,
		UpdatedAt: newEvent.UpdatedAt,
	}

	return response, nil
}

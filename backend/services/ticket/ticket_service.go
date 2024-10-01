package ticketService

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/codepnw/ticket-api/cmd/database"
	"github.com/codepnw/ticket-api/pkg/errs"
	eventRepository "github.com/codepnw/ticket-api/repositories/event"
	ticketRepository "github.com/codepnw/ticket-api/repositories/ticket"
	eventService "github.com/codepnw/ticket-api/services/event"
)

const timeLayout = "2006-01-02 15:04"

type ticketService struct {
	ticketRepo ticketRepository.ITicketRepository
}

func NewTicketService(ticketRepo ticketRepository.ITicketRepository) ITicketService {
	return &ticketService{ticketRepo: ticketRepo}
}

func (s *ticketService) GetTickets() ([]*TicketResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	tickets, err := s.ticketRepo.GetMany(context)
	if err != nil {
		log.Println(err)
		return nil, errs.NewErrUnexpected()
	}

	responses := []*TicketResponse{}

	for _, t := range tickets {
		// Event Data
		event, err := eventData(int(t.EventID))
		if err != nil {
			log.Println(err)
			return nil, errs.NewErrUnexpected()
		}

		responses = append(responses, &TicketResponse{
			ID:      t.ID,
			EventID: t.EventID,
			Event: eventRepository.Event{
				ID:        event.ID,
				Name:      event.Name,
				Location:  event.Location,
				Date:      event.Date,
				CreatedAt: event.CreatedAt,
				UpdatedAt: event.UpdatedAt,
			},
			Entered:   t.Entered,
			CreatedAt: t.CreatedAt,
			UpdatedAt: t.UpdatedAt,
		})
	}
	return responses, nil
}

func (s *ticketService) GetTicket(ticketID uint) (*TicketResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	ticket, err := s.ticketRepo.GetOne(context, ticketID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("ticket_id not found")
		}
		log.Println(err)
		return nil, errs.NewErrUnexpected()
	}

	// Event Data
	event, err := eventData(int(ticket.EventID))
	if err != nil {
		log.Println(err)
		return nil, errs.NewErrUnexpected()
	}

	response := &TicketResponse{
		ID:      ticket.ID,
		EventID: ticket.EventID,
		Event: eventRepository.Event{
			ID:        event.ID,
			Name:      event.Name,
			Location:  event.Location,
			Date:      event.Date,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		},
		Entered:   ticket.Entered,
		CreatedAt: ticket.CreatedAt,
		UpdatedAt: ticket.UpdatedAt,
	}
	return response, nil
}

func (s *ticketService) CreateTicket(request TicketRequest) (*TicketResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	// Event Data
	event, err := eventData(int(request.EventID))
	if err != nil {
		log.Println(err)
		return nil, errs.NewErrUnexpected()
	}

	ticket := &ticketRepository.Ticket{
		EventID: request.EventID,
		Event: eventRepository.Event{
			ID:        event.ID,
			Name:      event.Name,
			Location:  event.Location,
			Date:      event.Date,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		},
	}

	t, err := s.ticketRepo.CreateOne(context, ticket)
	if err != nil {
		log.Println(err)
		return nil, errs.NewErrUnexpected()
	}

	response := &TicketResponse{
		ID:        t.ID,
		EventID:   t.EventID,
		Event:     t.Event,
		Entered:   t.Entered,
		CreatedAt: time.Now().Format(timeLayout),
		UpdatedAt: time.Now().Format(timeLayout),
	}

	return response, nil
}

func (s *ticketService) ValidateTicket(ticketID uint, validate *ticketRepository.ValidateTicket) error {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	data := &ticketRepository.ValidateTicket{
		Entered: true,
	}

	err := s.ticketRepo.UpdateOne(context, ticketID, data)
	if err != nil {
		log.Println(err)
		return errs.NewErrUnexpected()
	}
	return nil
}

func eventData(eventID int) (*eventService.EventResponse, error) {
	eventRepo := eventRepository.NewEventRepository(database.GetDB())
	eventSrv := eventService.NewEventService(eventRepo)

	event, err := eventSrv.GetEvent(strconv.Itoa(eventID))
	if err != nil {
		return nil, err
	}

	return event, nil
}

package ticketService

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/codepnw/ticket-api/pkg/errs"
	ticketRepository "github.com/codepnw/ticket-api/repositories/ticket"
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
		responses = append(responses, &TicketResponse{
			ID:        t.ID,
			EventID:   t.EventID,
			Event:     t.Event,
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

	response := &TicketResponse{
		ID:        ticket.ID,
		EventID:   ticket.EventID,
		Event:     ticket.Event,
		Entered:   ticket.Entered,
		CreatedAt: ticket.CreatedAt,
		UpdatedAt: ticket.UpdatedAt,
	}
	return response, nil
}

func (s *ticketService) CreateTicket(request TicketRequest) (*TicketResponse, error) {
	context, cancel := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
	defer cancel()

	ticket := &ticketRepository.Ticket{
		EventID: request.EventID,
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

func (s *ticketService) UpdateTicket(ticketID uint, validate *ticketRepository.ValidateTicket) error {
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

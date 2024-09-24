package handlers

import (
	"net/http"

	eventService "github.com/codepnw/ticket-api/services/event"
	"github.com/gin-gonic/gin"
)

type eventHandler struct {
	eventSrv eventService.IEventService
}

func NewEventHandler(eventSrv eventService.IEventService) *eventHandler {
	return &eventHandler{eventSrv: eventSrv}
}

func (h *eventHandler) GetEvents(c *gin.Context) {
	events, err := h.eventSrv.GetEvents()
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *eventHandler) GetEvent(c *gin.Context) {
	id := c.Param("event_id")

	event, err := h.eventSrv.GetEvent(id)
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *eventHandler) NewEvent(c *gin.Context) {
	var event eventService.NewEventRequest

	if err := c.ShouldBindJSON(&event); err != nil {
		errorBadRequest(c, err.Error())
		return
	}

	response, err := h.eventSrv.CreateEvent(event)
	if err != nil {
		errorInternalServer(c, err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

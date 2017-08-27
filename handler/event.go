package handler

import (
	"log"
	"net/http"
	"time"

	"oxylus/driver/particleio"
	"oxylus/event"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

type (
	// Date represents the incoming date to trigger
	Date struct {
		Year   int `json:"year"`
		Month  int `json:"month"`
		Day    int `json:"day"`
		Hour   int `json:"hour"`
		Minute int `json:"minute"`
		Second int `json:"second"`
	}
	// EventRequest represents the request a user makes to create an event
	EventRequest struct {
		Date         Date   `json:"date"`
		Action       string `json:"action"`
		Repeats      bool   `json:"repeats"`
		Driver       string `json:"driver"`
		DeviceID     string `json:"deviceId"`
		AccessToken  string `json:"accessToken"`
		TimeInterval string `json:"timeInterval"`
	}
	// EventResponse represents the response from the server after an event is created
	EventResponse struct {
		Event    *event.Event        `json:"event"`
		Commands []map[string]string `json:"commands"`
		Request  *EventRequest       `json:"request"`
	}
)

// CreateEvent --
func (h *Handler) CreateEvent(c echo.Context) error {
	id := c.Param("id")
	var request EventRequest
	var err error
	if err = c.Bind(request); err != nil {
		return err
	}
	e := event.New()
	if request.TimeInterval == "" {
		e.FinishAt = time.Date(request.Date.Year, time.Month(request.Date.Month), request.Date.Day, request.Date.Hour, request.Date.Minute, request.Date.Second, 0, time.Local)
		e.TimeInterval = time.Until(e.FinishAt)
	} else {
		if e.TimeInterval, err = time.ParseDuration(request.TimeInterval); err != nil {
			log.Println(err)
		}
	}
	e.Action = request.Action
	e.Driver = &particleio.ParticleIO{
		UUID:        uuid.NewV4().String(),
		DeviceID:    request.DeviceID,
		AccessToken: request.AccessToken,
	}
	e.Repeats = request.Repeats
	// response := new(EventResponse)

	h.EventRegistry.Add(id, e)
	// move the start to the registry
	// if err = h.EventRegistry.StartTimer(id, e.UUID.String()); err != nil {
	// 	log.Println(err)
	// }
	// response.Event = e
	// response.Request = request
	// response.Commands = e.Driver.Commands()
	return c.NoContent(http.StatusCreated)
}

// DeleteEvent will stop a timer and delete the event from the registry
func (h *Handler) DeleteEvent(c echo.Context) error {
	id := c.Param("id")

	h.EventRegistry.StopTimer(id)
	h.EventRegistry.Remove(id)
	return c.NoContent(http.StatusOK)
}

// GetUserEvents will return all events for the user
func (h *Handler) GetUserEvents(c echo.Context) error {
	id := c.Param("id")
	events := h.EventRegistry.GetAll(id)
	var results = make([]*event.Event, len(events))
	idx := 0
	for _, value := range events {
		results[idx] = value
		idx++
	}
	return c.JSONPretty(http.StatusOK, results, "  ")
}

// GetUserEvent will return one event for the user
func (h *Handler) GetUserEvent(c echo.Context) error {
	id := c.Param("id")
	event, err := h.EventRegistry.Get(id)
	if err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, event, "  ")
}

// GetAllEvents will return all events in the registry
func (h *Handler) GetAllEvents(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, h.EventRegistry.Registry, "  ")
}

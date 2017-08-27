package handler

import (
	"log"
	"net/http"
	"oxylus/poller"
	"time"

	"oxylus/driver/particleio"

	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

// PollerRequest represents the payload received to create or update a poller.
type PollerRequest struct {
	Action       string `json:"action"`
	Driver       string `json:"driver"`
	DeviceID     string `json:"deviceID"`
	AccessToken  string `json:"accessToken"`
	PollInterval string `json:"pollInterval"`
	IsPolling    bool   `json:"isPolling"`
}

// AddPoller creates a poller
func (h *Handler) AddPoller(c echo.Context) error {
	id := c.Param("id")
	request := &PollerRequest{}
	var err error
	if err = c.Bind(request); err != nil {
		return err
	}
	p := &poller.Poller{}
	p.UUID = uuid.NewV4().String()
	p.Action = request.Action
	p.Driver = &particleio.ParticleIO{
		UUID:        p.UUID,
		DeviceID:    request.DeviceID,
		AccessToken: request.AccessToken,
	}
	p.PollInterval, err = time.ParseDuration(request.PollInterval)
	if err != nil {
		return err
	}
	p.IsPolling = request.IsPolling
	p.User = id

	db := h.DB.Clone()
	defer db.Close()

	if err := db.DB("oxylus").C("pollers").Insert(&p); err != nil {
		return err
	}
	// if ispolling then send the poller to the registry
	// turn this into a channel
	h.PollerRegistry.Add(p.UUID, p)
	return c.NoContent(http.StatusCreated)
}

// UpdatePoller will toggle the state of the poller
func (h *Handler) UpdatePoller(c echo.Context) error {
	id := c.Param("id")
	request := &PollerRequest{}
	var err error
	if err = c.Bind(request); err != nil {
		return err
	}

	db := h.DB.Clone()
	defer db.Close()

	newInterval, err := time.ParseDuration(request.PollInterval)
	if err != nil {
		log.Println(err)
	}
	if err := db.DB("oxylus").C("pollers").Update(
		bson.M{"uuid": id},
		bson.M{"$set": bson.M{
			"ispolling":    request.IsPolling,
			"action":       request.Action,
			"pollinterval": newInterval}}); err != nil {
		return err
	}

	if request.IsPolling {
		var p poller.Poller
		if err := db.DB("oxylus").C("pollers").Find(bson.M{"uuid": id}).One(&p); err != nil {
			return err
		}
		h.PollerRegistry.Add(id, &p)
	} else {
		h.PollerRegistry.Remove(id)
	}

	return c.NoContent(http.StatusOK)
}

// DeletePoller --
func (h *Handler) DeletePoller(c echo.Context) error {
	id := c.Param("id")
	db := h.DB.Clone()
	defer db.Close()

	if err := db.DB("oxylus").C("pollers").Remove(bson.M{"uuid": id}); err != nil {
		return err
	}

	// remove from registry
	h.PollerRegistry.Remove(id)
	return c.NoContent(http.StatusOK)
}

// GetPollers returns all pollers attached to a user
func (h *Handler) GetPollers(c echo.Context) error {
	id := c.Param("id")

	db := h.DB.Clone()
	defer db.Close()

	var p []*poller.Poller
	if err := db.DB("oxylus").C("pollers").Find(bson.M{"user": id}).All(&p); err != nil {
		return err
	}
	return c.JSONPretty(http.StatusOK, p, "  ")
}

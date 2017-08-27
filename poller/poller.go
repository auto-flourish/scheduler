package poller

import (
	"oxylus/driver/particleio"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Poller represents an object that repeats a task every n seconds
// This represents a contract with a hardware target
type Poller struct {
	ID           bson.ObjectId          `json:"id" bson:"_id,omitempty"`
	UUID         string                 `json:"uuid"`
	Action       string                 `json:"action"`
	Timer        *time.Timer            `json:"-"`
	PollInterval time.Duration          `json:"pollInterval"`
	Driver       *particleio.ParticleIO `json:"driver"`
	IsPolling    bool                   `json:"isPolling"`
	User         string                 `json:"user"`
}

// Poll is the loop that polls the hardware
func (p *Poller) Poll() (interface{}, error) {
	response, err := p.Driver.Poll(p.Action)
	if err != nil {
		return nil, err
	}
	return response, nil
}

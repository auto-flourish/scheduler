package event

import (
	"errors"
	"fmt"
	"oxylus/driver/particleio"
	"time"

	uuid "github.com/satori/go.uuid"
)

var (
	// ErrTimeIntervalNotPositive means the timer can't activate a function in the future
	ErrTimeIntervalNotPositive = errors.New("time interval is less than zero")
)

// Event represents an action to be executed at some time
type Event struct {
	UUID         uuid.UUID              `json:"uuid"`
	CreatedAt    time.Time              `json:"createdAt"`
	FinishAt     time.Time              `json:"finishAt"`
	Action       string                 `json:"action"`
	Driver       *particleio.ParticleIO `json:"driver"`
	Repeats      bool                   `json:"repeats"`
	TimeInterval time.Duration          `json:"timeInterval"`
	Timer        *time.Timer            `json:"-"`
	Status       string                 `json:"status"`
}

// StopTimer attempts to stop the time and returns a bool indicating status
func (e *Event) StopTimer() bool {
	return e.Timer.Stop()
}

// String returns this object as a string
func (e *Event) String() string {
	return fmt.Sprintf("%s %s %t %s", e.UUID.String(), e.Action, e.Repeats, e.FinishAt)
}

// New returns a new event
func New() *Event {
	return &Event{
		UUID:      uuid.NewV4(),
		CreatedAt: time.Now(),
		Status:    "created",
	}
}

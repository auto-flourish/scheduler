package eventregistry

import (
	"fmt"
	"oxylus/driver/particleio"
	"oxylus/event"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Ensure we can create a registry, create an event, add the event to the registry
// start the timer through the registry
// stop the timer through the registry
// we get events through out timerstart and timerended channels
func TestNewRegistry(t *testing.T) {
	registry := New()
	registryUser := uuid.NewV5(uuid.NamespaceURL, "registryUser")
	elementID := uuid.NewV5(uuid.NamespaceURL, "elementID")
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+3, 0, time.Local)
	driver := particleio.NewDriver()
	e := &event.Event{
		UUID:         elementID,
		FinishAt:     date,
		Driver:       driver,
		Action:       "test",
		CreatedAt:    time.Now(),
		Repeats:      false,
		TimeInterval: time.Until(date),
	}
	go func(r *EventRegistry) {
		for {
			select {
			case msg := <-r.TimerStarted:
				fmt.Println(msg.UUID.String() + " started")
			case msg := <-r.TimerEnded:
				fmt.Println(msg.UUID.String() + " ended")
			}
		}
	}(registry)
	registry.Add(registryUser.String(), e)
	registry.StartTimer(registryUser.String(), e.UUID.String())
	time.Sleep(time.Second * 5)
	registry.StopTimer(registryUser.String(), e.UUID.String())
}

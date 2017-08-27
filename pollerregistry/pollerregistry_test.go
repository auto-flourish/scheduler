package pollerregistry

import (
	"fmt"
	"oxylus/driver/particleio"
	"oxylus/poller"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestAddPoller(t *testing.T) {
	reg := New()
	dur, err := time.ParseDuration("3s")
	if err != nil {
		t.Fatal("could not parse duration")
	}
	p := &poller.Poller{
		UUID:   uuid.NewV4().String(),
		Action: "AllSensData",
		Driver: &particleio.ParticleIO{
			UUID:        uuid.NewV4(),
			DeviceID:    "3c0026000247353137323334",
			AccessToken: "f235b0985b6b46d0b50e3ee93e051dfe1742b201",
		},
		PollInterval: dur,
		IsPolling:    true,
		User:         "test-uuid",
	}
	go func(reg *PollerRegistry) {
		for {
			select {
			case msg := <-reg.ToDB:
				fmt.Println(msg)
				if msg == nil {
					t.Fatal("msg should not be nil")
				}
			}
		}
	}(reg)
	p.UUID = uuid.NewV4().String()
	reg.Add(p.UUID, p)
	time.Sleep(time.Second * 10)
}

func TestDeletePoller(t *testing.T) {
	reg := New()
	dur, err := time.ParseDuration("3s")
	if err != nil {
		t.Fatal("could not parse duration")
	}
	p := &poller.Poller{
		UUID:   uuid.NewV4().String(),
		Action: "AllSensData",
		Driver: &particleio.ParticleIO{
			UUID:        uuid.NewV4(),
			DeviceID:    "3c0026000247353137323334",
			AccessToken: "f235b0985b6b46d0b50e3ee93e051dfe1742b201",
		},
		PollInterval: dur,
		IsPolling:    true,
		User:         "test-uuid",
	}
	go func(reg *PollerRegistry) {
		for {
			select {
			case msg := <-reg.ToDB:
				fmt.Println(msg)
				if msg == nil {
					t.Fatal("msg should not be nil")
				}
			}
		}
	}(reg)
	p.UUID = uuid.NewV4().String()
	reg.Add(p.UUID, p)
	time.Sleep(time.Second * 5)
	reg.Remove(p.UUID)
}
func TestUpdatePoller(t *testing.T) {
	reg := New()
	dur, err := time.ParseDuration("3s")
	if err != nil {
		t.Fatal("could not parse duration")
	}
	p := &poller.Poller{
		UUID:   uuid.NewV4().String(),
		Action: "AllSensData",
		Driver: &particleio.ParticleIO{
			UUID:        uuid.NewV4(),
			DeviceID:    "3c0026000247353137323334",
			AccessToken: "f235b0985b6b46d0b50e3ee93e051dfe1742b201",
		},
		PollInterval: dur,
		IsPolling:    true,
		User:         "test-uuid",
	}
	go func(reg *PollerRegistry) {
		for {
			select {
			case msg := <-reg.ToDB:
				fmt.Println(msg)
				if msg == nil {
					t.Fatal("msg should not be nil")
				}
			}
		}
	}(reg)
	p.UUID = uuid.NewV4().String()
	reg.Add(p.UUID, p)
	time.Sleep(time.Second * 5)
	reg.Remove(p.UUID)
	newDur, err := time.ParseDuration("10s")
	if err != nil {
		t.Fatal("could not create new duration")
	}
	p.PollInterval = newDur
	reg.Add(p.UUID, p)
	time.Sleep(time.Second * 12)
}

// Ensure we can create a registry, create an event, add the event to the registry
// start the timer through the registry
// stop the timer through the registry
// we get events through out timerstart and timerended channels
// func TestNewRegistry(t *testing.T) {
// 	registry := New()
// 	registryUser := uuid.NewV5(uuid.NamespaceURL, "registryUser")
// 	elementID := uuid.NewV5(uuid.NamespaceURL, "elementID")
// 	now := time.Now()
// 	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+3, 0, time.Local)
// 	driver := particleio.NewDriver()
// 	e := &event.Event{
// 		UUID:         elementID,
// 		FinishAt:     date,
// 		Driver:       driver,
// 		Action:       "test",
// 		CreatedAt:    time.Now(),
// 		Repeats:      false,
// 		TimeInterval: time.Until(date),
// 	}
// 	go func(r *EventRegistry) {
// 		for {
// 			select {
// 			case msg := <-r.TimerStarted:
// 				fmt.Println(msg.UUID.String() + " started")
// 			case msg := <-r.TimerEnded:
// 				fmt.Println(msg.UUID.String() + " ended")
// 			}
// 		}
// 	}(registry)
// 	registry.Add(registryUser.String(), e)
// 	registry.StartTimer(registryUser.String(), e.UUID.String())
// 	time.Sleep(time.Second * 5)
// 	registry.StopTimer(registryUser.String(), e.UUID.String())
// }

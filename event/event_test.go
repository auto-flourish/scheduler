package event

import (
	"fmt"
	"oxylus/driver/particleio"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func TestTimerRunsNormally(t *testing.T) {
	// Test one shot timer runs normally
	events := make([]*Event, 1)
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+1, 0, time.Local)
	driver := particleio.New()
	e := &Event{
		UUID:         uuid.NewV4(),
		FinishAt:     date,
		Driver:       driver,
		Action:       "test",
		CreatedAt:    time.Now(),
		Repeats:      false,
		TimeInterval: time.Until(date),
	}
	events[0] = e
	if err := events[0].StartTimer(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(2 * time.Second)
	if events[0].timer.Stop() == true {
		t.Error("timer should have expired but isn't")
	}
	events[0] = nil
}

func TestTimerStopsNormally(t *testing.T) {
	// Test invalidation of a timer
	events := make([]*Event, 1)
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+2, 0, time.Local)
	driver := particleio.NewDriver()
	e := &Event{
		UUID:         uuid.NewV4(),
		FinishAt:     date,
		Driver:       driver,
		Action:       "test",
		CreatedAt:    time.Now(),
		User:         uuid.NewV5(uuid.NamespaceURL, "user-uuid"),
		Repeats:      false,
		TimeInterval: time.Until(date),
	}
	events[0] = e
	if err := events[0].StartTimer(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(1 * time.Second)
	if events[0].timer.Stop() == false {
		t.Error("timer shouldn't be expired but it is")
	}
	events[0] = nil
}

func TestTimerReapeats(t *testing.T) {
	// Test timer repeats normally
	events := make([]*Event, 1)
	now := time.Now()
	date := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+2, 0, time.Local)
	driver := particleio.NewDriver()
	e := &Event{
		UUID:         uuid.NewV4(),
		FinishAt:     date,
		Driver:       driver,
		Action:       "test",
		CreatedAt:    time.Now(),
		User:         uuid.NewV5(uuid.NamespaceURL, "user-uuid"),
		Repeats:      true,
		TimeInterval: time.Until(date),
	}
	events[0] = e
	if err := events[0].StartTimer(); err != nil {
		fmt.Println(err)
	}
	time.Sleep(10 * time.Second)
	if events[0].timer.Stop() == false {
		t.Error("timer shouldn't be expired but it is")
	}
	events[0] = nil
}

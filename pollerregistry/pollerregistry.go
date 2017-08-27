package pollerregistry

import (
	"errors"
	"log"
	"oxylus/poller"
	"sync"
	"time"
)

var (
	// ErrKeyNotFound denotes the key does not exist in the registry
	ErrKeyNotFound = errors.New("key not found")
	// ErrElementNotFound means the event id does not exist
	ErrElementNotFound = errors.New("element not found in registry")
	// ErrTimeIntervalLessThanZero means the time cant be scheduled
	ErrTimeIntervalLessThanZero = errors.New("time intveral is less than zero")
)

// PollerRegistry uses a uuid to key a list of events
type PollerRegistry struct {
	Registry map[string]*poller.Poller
	sync.Mutex
	ToDB         chan interface{}
	UpdateStatus chan string
}

// New returns a new registry
func New() *PollerRegistry {
	return &PollerRegistry{
		Registry:     make(map[string]*poller.Poller),
		ToDB:         make(chan interface{}),
		UpdateStatus: make(chan string),
	}
}

// Poll --
func (e *PollerRegistry) Poll(key string) error {
	poller, err := e.Get(key)
	if err != nil {
		return ErrElementNotFound
	}
	poller.Timer = time.AfterFunc(poller.PollInterval, func() {
		response, err := poller.Poll()
		if err != nil {
			log.Println(err)
			e.UpdateStatus <- poller.UUID
		} else {
			e.Poll(key)
			e.ToDB <- response
		}
	})
	return nil
}

// StopPoller allows us to place middleware around the event
// This is useful for pruning a dead event
func (e *PollerRegistry) StopPoller(key string) error {
	poller, err := e.Get(key)
	if err != nil {
		return ErrElementNotFound
	}
	poller.Timer.Stop()
	return nil
}

// GetAll returns the map against the user uuid
func (e *PollerRegistry) GetAll() map[string]*poller.Poller {
	return map[string]*poller.Poller{}
}

// Get returns a single event
func (e *PollerRegistry) Get(key string) (*poller.Poller, error) {
	if val, ok := e.Registry[key]; ok {
		return val, nil
	}
	return nil, ErrElementNotFound
}

// Add sets a value in the registry
func (e *PollerRegistry) Add(key string, val *poller.Poller) {
	if err := e.Remove(key); err != nil {
		log.Println(err)
	}

	e.Lock()
	e.Registry[key] = val
	e.Unlock()
	if err := e.Poll(key); err != nil {
		log.Println(err)
	}
}

// Remove will remove an event from the map
func (e *PollerRegistry) Remove(key string) error {
	if err := e.StopPoller(key); err != nil {
		return err
	}
	if _, ok := e.Registry[key]; !ok {
		e.Lock()
		delete(e.Registry, key)
		e.Unlock()
	}
	return nil
}

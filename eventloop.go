package event

import (
	"time"
)

// EventLoop is a struct that can contain many events
type EventLoop struct {
	// All events, protected by a mutex whenever it is used
	events []*Event
}

func NewEventLoop() *EventLoop {
	return &EventLoop{[]*Event{}}
}
func (s *EventLoop) AddEvent(e *Event) {
	s.events = append(s.events, e)
}

func (s *EventLoop) Go(sleep time.Duration) {
	// Use an endless event loop
	for {
		// For each possible event
		for _, e := range s.events {
			// Check if the event should trigger
			if e.ShouldTrigger() {
				// When triggering an event, run it in the background
				go e.Trigger()
			}
		}
		// How long to sleep before checking again
		time.Sleep(sleep)
	}
}

package event

import (
	"fmt"
	"sync"
	"time"
)

// Event stores a time window for when the event can be triggered (`from` up to `upTo`),
// how long it should take to cooldown before being able to be re-triggered, which
// action should be performed when triggered, when it was last triggered, a mutex
// and a boolean variable for keeping track of if the action is still ongoing or not.
type Event struct {
	from       time.Time
	upTo       time.Time
	cooldown   time.Duration // How long to cool down before retriggering
	actionFunc func()        // Action takes no arguments
	triggered  time.Time     // When was the event last triggered
	mutex      *sync.RWMutex
	ongoing    bool
	dateEvent  bool // Is the triggering based on the clock or date?
}

// New creates a new Event, that should happen at the given "when" time,
// within the given time window, with an associated cooldown period after the
// event has been triggered. The event can be retriggered after every cooldown,
// within the time window.
func New(when time.Time, window, cooldown time.Duration, action func()) *Event {
	e := &Event{when, when.Add(window), cooldown, action, time.Time{}, &sync.RWMutex{}, false, false}
	e.update()
	return e
}

func NewDateEvent(when time.Time, window, cooldown time.Duration, action func()) *Event {
	return &Event{when, when.Add(window), cooldown, action, time.Time{}, &sync.RWMutex{}, false, true}
}

// From is the time from when the event should be able to be triggered.
func (e *Event) From() time.Time {
	return e.from
}

// UpTo is the time where the event should no longer be able to be triggered.
func (e *Event) UpTo() time.Time {
	return e.upTo
}

// Cooldown is how long to wait after the event has been triggered, before
// being possible to trigger again.
func (e *Event) Cooldown() time.Duration {
	return e.cooldown
}

// Duration is for how long the window that this event can be triggered is
func (e *Event) Duration() time.Duration {
	return e.upTo.Sub(e.from)
}

// Between returns true if the given time t is between the two timestamps
// a (inclusive) and b (exclusive)
func Between(t, a, b time.Time) bool {
	return (t.Sub(a) >= 0) && (t.Sub(b) < 0)
}

// Update every field, except the hour/minute/second, to the current day.
func (e *Event) update() {
	now := time.Now()

	// Get hour, minute and second from the event
	hour, min, sec := e.from.Clock()
	// Schedule a new time
	e.from = time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, now.Nanosecond(), now.Location())

	// Get hour, minute and second from the event
	hour, min, sec = e.upTo.Clock()
	// Schedule a new time
	e.upTo = time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, now.Nanosecond(), now.Location())
}

// Has checks if the Event has time t in its interval:
// from p.From() and up to but not including p.UpTo()
func (e *Event) Has(t time.Time) bool {
	// If it's not a date-event, update the year/month/day to the current
	// year/month/day before checking
	if !e.dateEvent {
		e.update()
	}
	return Between(t, e.From(), e.UpTo())
}

// ShouldTrigger returns true if the current time is in the interval
// of the event AND it is not ongoing AND it is not in the cooldown period.
func (e *Event) ShouldTrigger() bool {
	t := time.Now()

	// Safely read the status
	e.mutex.RLock()
	retval := !e.ongoing && e.Has(t) && !Between(t, e.triggered, e.triggered.Add(e.cooldown))
	e.mutex.RUnlock()

	return retval
}

// Trigger triggers this event. The trigger time is noted, the associated
// action is performed and a cooldown period is initiated with time.Sleep.
// It is expected that this function will be called as a goroutine.
func (e *Event) Trigger() {
	// Safely update the status
	e.mutex.Lock()
	e.ongoing = true
	e.triggered = time.Now()
	e.mutex.Unlock()

	// Perform the action
	e.actionFunc()

	// If there is time left, sleep some
	passed := time.Now().Sub(e.triggered)
	time.Sleep(e.cooldown - passed)

	// Safely update the status
	e.mutex.Lock()
	e.ongoing = false
	e.mutex.Unlock()
}

// String returns a string with information about this event
func (e *Event) String() string {
	return fmt.Sprintf("Event from %s up to %s. Cooldown: %v. Should trigger: %v", e.from.Format("15:04:05"), e.upTo.Format("15:04:05"), e.cooldown, e.ShouldTrigger())
}

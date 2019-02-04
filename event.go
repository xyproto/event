package event

import (
	"fmt"
	"sync"
	"time"
)

// Event happens at a specific time, then is deactivated until the time interval has passed.
// It shoud only happen once per time interval, so ShouldTrigger will be set to false
// once it is done. All events that are not in the interval will have their ShouldTrigger
// status re-enabled.
// Event is a Triggerer.
type Event struct {
	from       time.Time
	upTo       time.Time
	actionFunc func()        // Action takes no arguments
	cooldown   time.Duration // how long to cool down before retriggering
	triggered  time.Time     // when was the event last triggered
	ongoing    bool
	mutex      *sync.RWMutex
}

// NewEvent creates a new Event, that should happen at the given "when" time, then retrigger after every cooldown, within the time window
func NewEvent(when time.Time, cooldown time.Duration, window time.Duration, action func()) *Event {
	return &Event{when, when.Add(window), action, cooldown, time.Time{}, false, &sync.RWMutex{}}
}

func (e *Event) Trigger() {
	e.mutex.Lock()
	e.ongoing = true
	e.triggered = time.Now()
	e.actionFunc()
	// If there is time left, sleep some
	passed := time.Now().Sub(e.triggered)
	time.Sleep(e.cooldown - passed)
	// Ok
	e.ongoing = false
	e.mutex.Unlock()
}

func (e *Event) From() time.Time {
	return e.from
}

func (e *Event) UpTo() time.Time {
	return e.upTo
}

func (e *Event) Cooldown() time.Duration {
	return e.cooldown
}

func (e *Event) Duration() time.Duration {
	return e.upTo.Sub(e.from)
}

// between checks if the given time t lies between the two timestamps
// a (inclusive) and b (exclusive)
func between(t, a, b time.Time) bool {
	return (t.Sub(a) >= 0) && (t.Sub(b) < 0)
}

func (e *Event) ShouldTrigger() bool {
	// If the time is in the allowed range AND it is not in the the time
	// between triggered and triggered+cooldown, it is active (and possible to
	// trigger)
	e.mutex.RLock()
	t := time.Now()
	retval := !e.ongoing && e.Has(t) && !between(t, e.triggered, e.triggered.Add(e.cooldown))
	e.mutex.RUnlock()
	return retval
}

func (e *Event) String() string {
	return fmt.Sprintf("Event from %s upto %s cooldown %v should trigger %v", e.from.Format("15:04:05"), e.upTo.Format("15:04:05"), e.cooldown, e.ShouldTrigger())
}

// Check if the Triggerer has the given time t in its time interval (from p.From() up to but not including p.UpTo())
func (e *Event) Has(t time.Time) bool {
	return between(t, e.From(), e.UpTo())
}

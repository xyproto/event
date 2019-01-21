package event

import (
	"fmt"
	"time"
)

// Event happens at a specific time, then is deactivated until the time interval has passed.
// It shoud only happen once per time interval, so Active will be set to false
// once it is done. All events that are not in the interval will have their Active
// status re-enabled.
// Event is a Performer.
type Event struct {
	from       time.Time
	upTo       time.Time
	actionFunc func()        // Action takes no arguments
	cooldown   time.Duration // how long to cool down before retriggering
	triggered  time.Time     // when was the event last triggered
}

func (e *Event) Perform() {
	e.triggered = time.Now()
	e.actionFunc()
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

func (e *Event) Active() bool {
	// If the time is in the allowed range AND it is not in the the time
	// between triggered and triggered+cooldown, it is active (and possible to
	// trigger)
	t := time.Now()
	return e.Has(t) && !Between(t, e.triggered, e.triggered.Add(e.cooldown))
}

func (e *Event) String() string {
	return fmt.Sprintf("Event [%v:%v) cooldown %v active %v\n", e.from, e.upTo, e.cooldown, e.Active())
}

// Check if the Performer has the given time t in its time interval (from p.From() up to but not including p.UpTo())
func (e *Event) Has(t time.Time) bool {
	return Between(t, e.From(), e.UpTo())
}

// NewEvent creates a new Event, that should happen once at the given "when" time, and no later than the event duration
func NewEvent(when time.Time, window time.Duration, action func()) *Event {
	return &Event{when, when.Add(window), action, window, time.Time{}}
}

// NewReEvent creates a new Event, that should happen at the given "when" time, then retrigger after every cooldown, within the time window
func NewReEvent(when time.Time, cooldown time.Duration, window time.Duration, action func()) *Event {
	return &Event{when, when.Add(window), action, cooldown, time.Time{}}
}

func (e *Event) SetActionProgressFunction(actionProgress func(float64)) {
	// Wrap the given function in a function that can measure the rate of progress
	e.actionFunc = func() {
		passed := time.Now().Sub(e.from) // how much time has passed?
		dur := e.upTo.Sub(e.from)        // how long should this event last?
		// Call the wrapped function, with an appropriate ratio
		var ratio float64
		if float64(dur) <= 0 {
			ratio = 0
		} else if float64(dur) >= 1.0 {
			ratio = 1
		} else {
			ratio = float64(passed) / float64(dur)
		}
		actionProgress(ratio)
	}
}

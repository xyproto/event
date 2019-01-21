package event

import (
	"fmt"
	"sync"
	"time"
)

type Performer interface { // What can an event do? It can perform actions when triggered, hence "Performer" (ref. io.Writer that can Write)
	From() time.Time                    // Trigger this event on or after this time (inclusive)
	UpTo() time.Time                    // Don't trigger the event on or after this time (exclusive)
	Cooldown() time.Duration            // How long to wait until the event can be repeated
	Perform()                           // Call the .action function in the struct. A time or progress argument is not needed, since the function itself can use time.Now() to find out.
	Duration() time.Duration            // calculated with the UpTo time minus the From time
	Has(time.Time) bool                 // check if a given point in time is within the Performer interval
	CooldownCounter() time.Duration     // How long is left of this cooldown until the next retrigger
	CooldownCounterSub(t time.Duration) // Subtract this duration from the cooldown countdown
	ResetCooldownCounter()              // Set the cooldown counter to the cooldown amount of time
	String() string
}

// Event happens at a specific time, then is deactivated until the time interval has passed.
// It shoud only happen once per time interval, so Active will be set to false
// once it is done. All events that are not in the interval will have their Active
// status re-enabled.
// Event is a Performer.
type Event struct {
	from              time.Time
	upTo              time.Time
	actionFunc        func()        // Action takes no arguments
	cooldown          time.Duration // how long to cool down before retriggering
	cooldownCountdown time.Duration // decrease this one every time the event loop sleep
}

func (e *Event) Perform() {
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

// Return hop long is left of the cooldown until next retrigger
func (e *Event) CooldownCounter() time.Duration { // How long is left of this cooldown until the next retrigger
	return e.cooldownCountdown
}

// Subtract time from the cooldown countdown. Will not go below 0.
func (e *Event) CooldownCounterSub(t time.Duration) { // Subtract this duration from the cooldown countdown
	e.cooldownCountdown -= t
	if e.cooldownCountdown < 0 {
		e.cooldownCountdown = 0
	}
}

// Reset the cooldown countdown to a cooldown amount of time
func (e *Event) ResetCooldownCounter() {
	e.cooldownCountdown = e.cooldown
}

func (e *Event) String() string {
	return fmt.Sprintf("Event [%v:%v) cooldown %v cooldownCountdown %v\n", e.from, e.upTo, e.cooldown, e.cooldownCountdown)
}

// NewEvent creates a new Event, that should happen at the given "from" time,
func NewEvent(from time.Time, duration time.Duration, action func()) *Event {
	return &Event{from, from.Add(duration), action, duration, duration}
}

// NewReEvent creates a new Event, that should happen at the given "from" time, then retrigger after every cooldown, within the time interval
func NewReEvent(from time.Time, duration time.Duration, action func(), cooldown time.Duration) *Event {
	return &Event{from, from.Add(duration), action, cooldown, cooldown}
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

// Check if the given time t lies between the two timestamps a (inclusive) and b (exclusive)
func Between(t, a, b time.Time) bool {
	return (t.Sub(a) >= 0) && (t.Sub(b) < 0)
}

// Check if the Performer has the given time t in its time interval (from p.From() up to but not including p.UpTo())
func (e *Event) Has(t time.Time) bool {
	return Between(t, e.From(), e.UpTo())
}

// All events, protected by a mutex whenever it is used
var performers = []Performer{}
var performerMutex = &sync.RWMutex{}

func AddEvent(e *Event) {
	p := Performer(e) // convert *Event to the Performer interface type
	performerMutex.Lock()
	performers = append(performers, p)
	performerMutex.Unlock()

}

func AddPerformer(p Performer) {
	performerMutex.Lock()
	performers = append(performers, p)
	performerMutex.Unlock()
}

func EventLoop() {
	// Use a single endless event loop
	for {
		fmt.Println("TIME", time.Now().String())

		// Events
		performerMutex.RLock()

		// Initial smallest cooldown value
		smallestCooldown := 200 * time.Millisecond
		if len(performers) > 0 {
			smallestCooldown = performers[0].Cooldown()
		}
		for _, performer := range performers {
			if performer.Cooldown() < smallestCooldown {
				smallestCooldown = performer.Cooldown()
			}
			if performer.Has(time.Now()) && (performer.CooldownCounter() == 0) {
				// Output info about the performer
				fmt.Println(performer)
				// Run the action in the background, and disable events that are in the time interval
				go performer.Perform()
				// Reset the cooldown counter
				performer.ResetCooldownCounter()
			} else {
				// Register how long this performer is going to sleep for
				performer.CooldownCounterSub(smallestCooldown)
			}
		}
		performerMutex.RUnlock()

		// Sleep, but no longer than the smallest cooldown
		time.Sleep(smallestCooldown)
	}
}

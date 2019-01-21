package event

import (
	"fmt"
	"sync"
	"time"
)

type Performer interface { // What can an event do? It can perform actions when triggered, hence "Performer" (ref. io.Writer that can Write)
	From() time.Time         // Trigger this event on or after this time (inclusive)
	UpTo() time.Time         // Don't trigger the event on or after this time (exclusive)
	Cooldown() time.Duration // How long to wait until the event can be repeated
	Perform()                // Call the .action function in the struct. A time or progress argument is not needed, since the function itself can use time.Now() to find out.
	Duration() time.Duration // calculated with the UpTo time minus the From time
	Has(time.Time) bool      // check if a given point in time is within the Performer interval
	Active() bool            // is this even active right now? (in the correct time interval, and not in the cooldown period)
	String() string          // string representation
}

func (s *Stage) AddPerformer(p Performer) {
	s.mut.Lock()
	s.performers = append(s.performers, p)
	s.mut.Unlock()
}

// Stage is a struct that can contain many performers
type Stage struct {
	// All events, protected by a mutex whenever it is used
	performers []Performer
	mut        *sync.RWMutex
}

func NewStage() *Stage {
	return &Stage{[]Performer{}, &sync.RWMutex{}}
}

func (s *Stage) AddEvent(e *Event) {
	p := Performer(e) // convert *Event to the Performer interface type
	s.mut.Lock()
	s.performers = append(s.performers, p)
	s.mut.Unlock()
}

func (s *Stage) EventLoop() {
	// Use a single endless event loop
	for {
		fmt.Println("TIME", time.Now().String())

		// Events
		s.mut.RLock()

		// Initial smallest cooldown value
		smallestCooldown := 500 * time.Millisecond
		if len(s.performers) > 0 {
			smallestCooldown = s.performers[0].Cooldown()
		} else {
			fmt.Println("no events")
		}
		for _, performer := range s.performers {
			if performer.Cooldown() < smallestCooldown {
				smallestCooldown = performer.Cooldown()
			}
			if performer.Active() {
				// Output info about the performer
				fmt.Println(performer)
				// Run the action in the background, and disable events that are in the time interval (this also records the trigger time)
				go performer.Perform()
			}
		}
		s.mut.RUnlock()

		// Sleep, but no longer than the smallest cooldown
		time.Sleep(smallestCooldown)
	}
}

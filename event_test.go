package event

import (
	"fmt"
	"testing"
	"time"
)

func Examplebetween() {
	before := time.Now()
	A := time.Now()
	t := time.Now()
	B := time.Now()
	after := time.Now()

	fmt.Println(between(t, A, B))          // true
	fmt.Println(between(t, before, A))     // false
	fmt.Println(between(t, B, after))      // false
	fmt.Println(between(t, before, after)) // true
	fmt.Println(between(t, before, B))     // true
	fmt.Println(between(t, A, after))      // true
	fmt.Println(between(t, after, before)) // false

	fmt.Println(between(t, t, B)) // true (from inclusive, to exclusive)
	fmt.Println(between(t, A, t)) // false (from inclusive, to exclusive)

	// Output:
	// true
	// false
	// false
	// true
	// true
	// true
	// false
	// true
	// false
}

func addEvents(s *EventLoop) {
	in0sec := time.Now()

	in2sec := in0sec.Add(2 * time.Second)
	in5sec := in0sec.Add(5 * time.Second)
	in15sec := in0sec.Add(15 * time.Second)
	in22sec := in0sec.Add(22 * time.Second)

	// Create two new events
	s.AddEvent(NewEvent(in5sec, 1*time.Second, 5*time.Second, func() {
		fmt.Println("This happens after 5s, cooldown: 1s, window: 5s")
	}))
	s.AddEvent(NewEvent(in15sec, 2*time.Second, 2*time.Second, func() {
		fmt.Println("This happens once after 15s")
	}))
	s.AddEvent(NewEvent(in2sec, 3*time.Second, 30*time.Second, func() {
		fmt.Println("This happens every 3s, within a 30 second window")
	}))
	s.AddEvent(NewEvent(in0sec, 2*time.Second, 20*time.Second, ProgressWrapper(in2sec, in22sec, func(p float64) {
		fmt.Printf("This happens every 2s, within a 20s window: %d%% complete\n", int(p*100))
	})))

}

//func addTransitions() {
//	in10sec := time.Now().Add(10 * time.Second)
//	in20sec := time.Now().Add(20 * time.Second)
//
//	// Create two new events
//	AddReEvent(NewTransition(in10sec, 2*time.Second, func(progress float64) {
//		fmt.Println("This event happens after 10 seconds, within a 2 second window! Progress:", progress)
//	}))
//	AddReEvent(NewTransition(in20sec, 2*time.Second, func(progress float64) {
//		fmt.Println("This event happens after 20 seconds, within a 2 second window! Progress:", progress)
//	}))
//}

func TestEventLoop(t *testing.T) {
	s := NewEventLoop()
	addEvents(s)
	//addTransitions()

	// TODO: Find the smallest cooldown and use that as the event loop delay
	fmt.Println("Run the event loop for 40s, with a loop delay of 50ms")
	go s.Go(50 * time.Millisecond)

	time.Sleep(40 * time.Second)
}

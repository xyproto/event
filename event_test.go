package event

import (
	"fmt"
	"testing"
	"time"
)

func ExampleBetween() {
	before := time.Now()
	A := time.Now()
	t := time.Now()
	B := time.Now()
	after := time.Now()

	fmt.Println(Between(t, A, B))          // true
	fmt.Println(Between(t, before, A))     // false
	fmt.Println(Between(t, B, after))      // false
	fmt.Println(Between(t, before, after)) // true
	fmt.Println(Between(t, before, B))     // true
	fmt.Println(Between(t, A, after))      // true
	fmt.Println(Between(t, after, before)) // false

	fmt.Println(Between(t, t, B)) // true (from inclusive, to exclusive)
	fmt.Println(Between(t, A, t)) // false (from inclusive, to exclusive)

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

func addEvents() {
	in5sec := time.Now().Add(5 * time.Second)
	in15sec := time.Now().Add(15 * time.Second)

	// Create two new events
	AddEvent(NewEvent(in5sec, 2*time.Second, func() {
		fmt.Println("This event happens after 5 seconds, within a 2 second window")
	}))
	AddEvent(NewEvent(in15sec, 200*time.Millisecond, func() {
		fmt.Println("This event happens after 15 seconds, within a 200 millisecond window")
	}))
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
	addEvents()
	//addTransitions()
	fmt.Println("Running the event loop for 40 seconds")
	go EventLoop()
	time.Sleep(40 * time.Second)
}
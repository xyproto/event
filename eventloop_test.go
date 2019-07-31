package event

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func fillEvents() *Loop {
	events := NewLoop()

	in0sec := time.Now()

	in01sec := in0sec.Add(100 * time.Millisecond)
	in2sec := in0sec.Add(2 * time.Second)
	in5sec := in0sec.Add(5 * time.Second)
	in15sec := in0sec.Add(15 * time.Second)
	in22sec := in0sec.Add(22 * time.Second)

	// Create four new events
	events.Add(New(in5sec, 5*time.Second, 1*time.Second, func() {
		fmt.Println("A happens after 5s, cooldown: 1s, window: 5s")
	}))
	events.Add(New(in15sec, 2*time.Second, 2*time.Second, func() {
		fmt.Println("B happens once after 15s")
	}))
	events.Add(New(in2sec, 30*time.Second, 3*time.Second, func() {
		fmt.Println("C happens every 3s, within a 30 second window")
	}))
	events.Add(New(in01sec, 20*time.Second, 2*time.Second, ProgressWrapperInterval(in2sec, in22sec, 2*time.Second, func(p float64) {
		fmt.Printf("D happens every 2s, within a 20s window: %d%% complete\n", int(p*100))
	})))

	// Add one-time events for the time markers
	fs := make([]func(), 5, 5)
	counter := 0
	for i := 0; i < 40; i += 10 {
		// Needed for the integer i to be enclosed correctly in the closure below
		seconds := i
		fs[counter] = func(seconds, i int) func() {
			return func() {
				fmt.Printf("--- %d second mark, i: %d ---\n", seconds, i)
			}
		}(seconds, i)
		f := fs[counter]
		var ptr uintptr = reflect.ValueOf(f).Pointer()
		fmt.Printf("registered function %T %v\n", f, ptr)
		events.Once(in0sec.Add(time.Duration(seconds)*time.Second), f)
		counter++
	}

	return events
}

func TestEventLoop(t *testing.T) {
	eventloop := fillEvents()

	// TODO: Find the smallest cooldown and use that as the event loop delay
	fmt.Println("Run the event loop for 30s, with a loop delay of 50ms")
	go eventloop.Go(50 * time.Millisecond)

	// Wait 30 seconds before returning
	time.Sleep(30 * time.Second)
}

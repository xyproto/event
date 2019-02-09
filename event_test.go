package event

import (
	"fmt"
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

func ExampleBetweenClock() {
	now := time.Now()
	in2secToday := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+2, now.Nanosecond(), now.Location())
	in2secNextYear := time.Date(now.Year()+1, now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()+2, now.Nanosecond(), now.Location())

	start := now

	//fmt.Println(in2secToday)
	//fmt.Println(in2secNextYear)

	in4sec := now.Add(4 * time.Second)

	fmt.Println(Between(in2secNextYear, start, in4sec))
	fmt.Println(Between(in2secToday, start, in4sec))
	fmt.Println(BetweenClock(in2secNextYear, start, in4sec))
	fmt.Println(BetweenClock(in2secToday, start, in4sec))

	// Output:
	// false
	// true
	// true
	// true
}

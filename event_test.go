package event

import (
	"fmt"
	"strings"
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
	fmt.Println(Between(in2secNextYear, start, in4sec))
	fmt.Println(Between(in2secToday, start, in4sec))

	// Output:
	// false
	// true
	// false
	// true
}

// dFmt formats a duration nicely
func dFmt(d time.Duration) string {
	s := fmt.Sprintf("%6s", d)
	if strings.Contains(s, ".") {
		pos := strings.Index(s, ".")
		if strings.Contains(s[pos:], "s") {
			s = s[:pos] + "s"
		}
	}
	if strings.HasSuffix(s, "m0s") {
		s = s[:len(s)-2]
	}
	if strings.HasSuffix(s, "h0m") {
		s = s[:len(s)-2]
	}
	return strings.TrimSpace(s)
}

func ExampleBetweenClockDay() {
	now := time.Now()

	noon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())
	two := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location())

	h1 := time.Hour * 1
	fmt.Println("one hour:", dFmt(h1))

	e := &Event{from: noon, upTo: two, cooldown: h1, clockOnly: true}
	fmt.Println("two hours:", dFmt(e.Duration()))

	// Output:
	// one hour: 1h
	// two hours: 1h59m59s
}

func ExampleBetweenClockDayBackwards() {
	now := time.Now()

	noon := time.Date(now.Year(), now.Month(), now.Day(), 14, 0, 0, 0, now.Location())
	two := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	h1 := time.Hour * 1
	fmt.Println("one hour:", dFmt(h1))

	e := &Event{from: noon, upTo: two, cooldown: h1, clockOnly: true}
	fmt.Println("twenty two hours:", dFmt(e.Duration()))

	// Output:
	// one hour: 1h
	// twenty two hours: 21h59m59s
}

func ExampleBetweenClockMidnight() {
	now := time.Now()

	// One hour before midnight, the date does not matter
	beforeMidnight := time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())

	// One hour after midnight, the date does not matter
	afterMidnight := time.Date(now.Year(), now.Month(), now.Day(), 1, 0, 0, 0, now.Location())

	h1 := time.Hour * 1
	fmt.Println("one hour:", dFmt(h1))

	e := &Event{from: beforeMidnight, upTo: afterMidnight, cooldown: h1, clockOnly: true}
	fmt.Println("two hours:", dFmt(e.Duration()))

	// Output:
	// one hour: 1h
	// two hours: 1h59m59s
}

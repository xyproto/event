package event

import (
	"time"
)

func ProgressWrapper(from, upto time.Time, progressFunction func(float64)) func() {
	start := time.Now()
	duration := upto.Sub(from)
	// Wrap the given function in a function that can measure the rate of progress
	return func() {
		start := start
		duration := duration
		passed := time.Now().Sub(start)
		ratio := 0.0
		if duration > 0 {
			ratio = float64(passed) / float64(duration)
		}
		// Clamp the ratio
		if ratio > 1.0 {
			ratio = 1.0
		}
		// Call the wrapped function, with an appropriate ratio
		progressFunction(ratio)
	}
}

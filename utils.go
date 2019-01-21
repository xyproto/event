package event

import (
	"time"
)

// Between checks if the given time t lies between the two timestamps a (inclusive) and b (exclusive)
func Between(t, a, b time.Time) bool {
	return (t.Sub(a) >= 0) && (t.Sub(b) < 0)
}

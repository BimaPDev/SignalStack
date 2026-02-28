package backoff

import (
	"math"
	"time"
)


func NextDelay(attemp int) time.Duration {
	baseDelay := 2
	maxDelay := 4
	calculate := int(math.Pow(2, float64(attemp)))
	delay := min(calculate*baseDelay, maxDelay)
	return time.Duration(delay) * time.Second
}

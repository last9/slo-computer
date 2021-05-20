package slo

import (
	"math"
	"time"
)

const (
	aDay   = 24 * time.Hour
	MinMtr = 1 * time.Hour
)

// convert duration to float seconds
func dToFS(d time.Duration) float64 {
	return float64(d.Seconds())
}

// convert float seconds to duration
func fsToD(f float64) time.Duration {
	return time.Duration(int64(f) * int64(time.Second))
}

// find minimum of given durations
func minD(d ...time.Duration) time.Duration {
	min := d[0]
	for _, right := range d[1:] {
		if min > right {
			min = right
		}
	}
	return min
}

// find maximum of given durations
func maxD(d ...time.Duration) time.Duration {
	max := d[0]
	for _, right := range d[1:] {
		if max < right {
			max = right
		}
	}
	return max
}

type Alerter interface {
	String() string
}

const float64EqualityThreshold = 0.005

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

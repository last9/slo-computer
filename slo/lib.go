package slo

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

type SLO struct {
	SLO        float64       // Desired SLO to achieve
	Throughput float64       // throughput (rpm)
	Period     time.Duration // SLO window
}

const (
	aDay    = time.Duration(24 * time.Hour)
	MIN_MTR = time.Duration(1 * time.Hour)
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

func NewSLO(period time.Duration, throughput, slo float64) (*SLO, error) {
	if period > 30*aDay {
		return nil, errors.New("period must <= 30 days")
	} else if period > aDay && period%aDay != 0 {
		return nil, errors.New("period must be a multiple of 24 hours")
	} else if period%3600 != 0 {
		return nil, errors.New("period must be a multiple of hours")
	} else if period <= 2*MIN_MTR {
		return nil, errors.New("period cannot be so low")
	}

	return &SLO{
		SLO:        slo,
		Throughput: throughput,
		Period:     period,
	}, nil
}

func (s *SLO) String() string {
	return fmt.Sprintf(
		"SLO of %.3f%% over %s for %.3f rpm", s.SLO, s.Period, s.Throughput,
	)
}

func (s *SLO) errorRate(errCount float64) float64 {
	return errCount / s.Throughput
}

func (s *SLO) burnRate(error_rate float64) float64 {
	return (error_rate * 100) / (100 - s.SLO)
}

func (s *SLO) budgetSpent(error_rate float64, alert_duration time.Duration) float64 {
	burn := s.burnRate(error_rate)
	return (burn * dToFS(alert_duration)) / dToFS(s.Period)
}

func (s *SLO) timeToExhaust(error_rate float64) time.Duration {
	burn := s.burnRate(error_rate)
	return fsToD(dToFS(s.Period) / burn)
}

func (s *SLO) errorImpact(errCount float64, duration time.Duration) *Impact {
	error_rate := s.errorRate(errCount)
	bs := s.budgetSpent(error_rate, duration)
	after := s.timeToExhaust(error_rate)
	return &Impact{errCount, error_rate, bs, duration, after}
}

type Impact struct {
	Errors      float64
	ErrorRate   float64
	BudgetSpent float64
	Duration    time.Duration
	BreaksAfter time.Duration
}

type AlertWindow struct {
	Name          string
	ErrorRate     float64
	BurnRate      float64
	BudgetSpent   float64
	TimeToExhaust time.Duration
	ShortWindow   time.Duration
	LongWindow    time.Duration
}

func (a *AlertWindow) String() string {
	return fmt.Sprintf(
		`
		Alert if error_rate > %.03f for last [%s] and also last [%s]
		This alert will trigger once %.2f%% of error budget is consumed,
		and leaves %s before the SLO is defeated.
		`,
		a.ErrorRate, a.LongWindow, a.ShortWindow,
		a.BudgetSpent*100, a.TimeToExhaust,
	)
}

func NewAlertWindow(
	slo *SLO, name string, error_rate float64, window time.Duration,
) *AlertWindow {
	a := AlertWindow{
		Name: name, ErrorRate: error_rate,
		ShortWindow: maxD(window/12, time.Duration(2*time.Minute)),
		LongWindow:  window,
	}

	a.BurnRate = slo.burnRate(error_rate)
	a.BudgetSpent = slo.budgetSpent(error_rate, window)
	a.TimeToExhaust = slo.timeToExhaust(error_rate)
	return &a
}

func AlertCalculator(s *SLO) []*AlertWindow {
	// Types of error-budget alerts
	out := make([]*AlertWindow, 2)

	// A good starting point for a fast-burn threshold policy is 10x the
	// baseline with a short (1- or 2-hour) lookback period.
	fastErrorRate := (100 - s.SLO) * 10 / 100

	// A good starting point for a slow-burn threshold is 2x the baseline with
	// a 24-hour lookback period.
	slowErrorRate := (100 - s.SLO) * 2 / 100

	var slowDuration time.Duration
	if s.Period > aDay { // SLO period is an order of a day
		slowDuration = aDay
	} else { // SLO period is an order of hours
		slowDuration = minD(time.Duration(2*time.Hour), s.Period/2)
	}

	fastDuration := maxD(slowDuration/24, time.Duration(5*time.Minute))

	// Slow-burn alert, which warns you of a rate of consumption that, if not
	// altered, exhausts your error budget before the end of the compliance
	// period. This type of condition is less urgent than a fast-burn
	// condition. "We are slightly exceeding where we'd like to be at this
	// point in the month, but we aren't in big trouble yet."
	// For a slow-burn alert, use a longer lookback period to smooth out
	// variations in shorter-term consumption.
	// The threshold you alert on in a slow-burn alert is higher than the ideal
	// performance for the lookback period, but not significantly higher. A
	// policy based on a shorter lookback period with high threshold might
	// generate too many alerts, even if the longer-term consumption levels
	// out. But if the consumption stays even a little too high for a longer
	// period, it eventually consumes all of your error budget.
	out[0] = NewAlertWindow(s, "slow", slowErrorRate, slowDuration)

	// When setting up alerting policies to monitor your error budget, it's a
	// good idea to set up two related alerting policies:
	// Fast-burn alert, which warns you of a sudden, large change in
	// consumption that, if uncorrected, will exhaust your error budget very
	// soon. "At this rate, we'll burn through the whole month's error budget
	// in two days!"
	// For a fast-burn alert, use a shorter lookback period so you are notified
	// quickly if a potentially disastrous condition has emerged and persisted,
	// even briefly. If it is truly disastrous, you don't want to wait long to
	// notice it.
	out[1] = NewAlertWindow(s, "fast", fastErrorRate, fastDuration)

	return out
}

// IsLowTraffic Finds the burn rate in 5 minutes
// If the throughput cannot withstand a small spike of 10/tps over 5 minutes
// probably the throughput and maturity is too low
func IsLowTraffic(slo *SLO) (*Impact, bool) {
	spikeImpact := slo.errorImpact(10.0, time.Duration(5*time.Minute))

	// If a single mini break is not fit enough to survive a MTTR. abort.
	if spikeImpact.BreaksAfter < 2*MIN_MTR {
		return spikeImpact, true
	}

	return spikeImpact, false
}

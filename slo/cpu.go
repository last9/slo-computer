package slo

import (
	"github.com/pkg/errors"
	"math"
	"time"
)

// Information about all instance types that support burst CPU performance
var InstanceTypes = map[string]EC2Instance{
	"t2.nano": {
		CreditRate:          3,
		MaxCredits:          72,
		VCPUs:               1,
		BaselineUtilization: 5,
	},
	"t2.micro": {
		CreditRate:          6,
		MaxCredits:          144,
		VCPUs:               1,
		BaselineUtilization: 10,
	},
	"t2.small": {
		CreditRate:          12,
		MaxCredits:          288,
		VCPUs:               1,
		BaselineUtilization: 20,
	},
	"t2.medium": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t2.large": {
		CreditRate:          36,
		MaxCredits:          864,
		VCPUs:               2,
		BaselineUtilization: 30,
	},
	"t2.xlarge": {
		CreditRate:          54,
		MaxCredits:          1296,
		VCPUs:               4,
		BaselineUtilization: 22.5,
	},
	"t2.2xlarge": {
		CreditRate:          81.6,
		MaxCredits:          1958.4,
		VCPUs:               8,
		BaselineUtilization: 17,
	},
	"t3.nano": {
		CreditRate:          6,
		MaxCredits:          144,
		VCPUs:               2,
		BaselineUtilization: 5,
	},
	"t3.micro": {
		CreditRate:          12,
		MaxCredits:          288,
		VCPUs:               2,
		BaselineUtilization: 10,
	},
	"t3.small": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t3.medium": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t3.large": {
		CreditRate:          36,
		MaxCredits:          864,
		VCPUs:               2,
		BaselineUtilization: 30,
	},
	"t3.xlarge": {
		CreditRate:          96,
		MaxCredits:          2304,
		VCPUs:               4,
		BaselineUtilization: 40,
	},
	"t3.2xlarge": {
		CreditRate:          192,
		MaxCredits:          4608,
		VCPUs:               8,
		BaselineUtilization: 40,
	},
	"t3a.nano": {
		CreditRate:          6,
		MaxCredits:          144,
		VCPUs:               2,
		BaselineUtilization: 5,
	},
	"t3a.micro": {
		CreditRate:          12,
		MaxCredits:          288,
		VCPUs:               2,
		BaselineUtilization: 10,
	},
	"t3a.small": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t3a.medium": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t3a.large": {
		CreditRate:          36,
		MaxCredits:          864,
		VCPUs:               2,
		BaselineUtilization: 30,
	},
	"t3a.xlarge": {
		CreditRate:          96,
		MaxCredits:          2304,
		VCPUs:               4,
		BaselineUtilization: 40,
	},
	"t3a.2xlarge": {
		CreditRate:          192,
		MaxCredits:          4608,
		VCPUs:               8,
		BaselineUtilization: 40,
	},
	"t4g.nano": {
		CreditRate:          6,
		MaxCredits:          144,
		VCPUs:               2,
		BaselineUtilization: 5,
	},
	"t4g.micro": {
		CreditRate:          12,
		MaxCredits:          288,
		VCPUs:               2,
		BaselineUtilization: 10,
	},
	"t4g.small": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t4g.medium": {
		CreditRate:          24,
		MaxCredits:          576,
		VCPUs:               2,
		BaselineUtilization: 20,
	},
	"t4g.large": {
		CreditRate:          36,
		MaxCredits:          864,
		VCPUs:               2,
		BaselineUtilization: 30,
	},
	"t4g.xlarge": {
		CreditRate:          96,
		MaxCredits:          2304,
		VCPUs:               4,
		BaselineUtilization: 40,
	},
	"t4g.2xlarge": {
		CreditRate:          192,
		MaxCredits:          4608,
		VCPUs:               8,
		BaselineUtilization: 40,
	},
}

type EC2Instance struct {
	CreditRate          float64
	MaxCredits          float64
	VCPUs               int
	BaselineUtilization float64
}

type BurstCPU struct {
	Instance          EC2Instance
	Utilisation float64
	Duration time.Duration
}

type BurstWindow struct {
	Name          string
	Utilisation     float64
	CreditBurnRate float64
	TimeToExhaust time.Duration
	ShortWindow   time.Duration
	LongWindow    time.Duration
}

func (b *BurstWindow) timeToExhaust(c *BurstCPU, utilisation float64, window time.Duration) (time.Duration, error) {
	credits := c.Instance.MaxCredits
	debitRate := b.creditsBurnRate(c, utilisation, window)

	if debitRate < c.Instance.CreditRate {
		return time.Duration(math.Inf(1)), errors.New("The instance will never run out of credits since it is underutilised.")
	}

	return  time.Duration((credits / debitRate) * dToFS(time.Hour)), nil
}

func (b *BurstWindow) creditsBurnRate(c *BurstCPU, utilisation float64, window time.Duration) float64 {
	creditsBurned := utilisation * float64(c.Instance.VCPUs) * (dToFS(window)/dToFS(time.Minute))
	return creditsBurned / (dToFS(window)/dToFS(time.Hour))
}

func NewBurstCPU(instance string, utilisation float64, duration time.Duration) (*BurstCPU, error){
	if _, ok := InstanceTypes[instance]; !ok {
		return nil, errors.New("Unknown instance type")
	}

	return &BurstCPU{
		Instance:    InstanceTypes[instance],
		Utilisation: utilisation,
		Duration: duration,

	}, nil
}

func NewBurstWindow(
	c *BurstCPU, name string, utilisation float64, window time.Duration,
) *BurstWindow {
	a := BurstWindow{
		Name: name,
		Utilisation: utilisation,
		ShortWindow: maxD(window/12, 2*time.Minute),
		LongWindow:  window,
	}

	a.CreditBurnRate = a.creditsBurnRate(c, utilisation, window)
	a.TimeToExhaust, err := a.timeToExhaust(c, utilisation, window)
	return &a
}

func BurstCalculator(b *BurstCPU) []*BurstWindow {
	// Types of CPU burst alerts
	out := make([]*BurstWindow, 2)

	// A short term burst in CPU can be characterized as double
	//the baseline utilisation for the instance
	fastUtilisationRate := b.Instance.BaselineUtilization * 2

	// A small increase in CPU usage can be characterized as the
	//usage being 20% above the baseline
	slowUtilisationRate := b.Instance.BaselineUtilization * 1.2

	var slowDuration time.Duration
	if b.Duration > aDay { // Duration is an order of a day
		slowDuration = aDay
	} else { // Duration is an order of hours
		slowDuration = minD(6*time.Hour, b.Duration/2)
	}

	fastDuration := maxD(slowDuration/24, 30*time.Minute)

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

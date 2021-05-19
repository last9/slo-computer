package slo

import (
	"embed"
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/pkg/errors"
)

//go:embed aws_instances.json
var instanceTypesFile embed.FS

// Information about all instance types that support burst CPU performance
var instanceTypes = map[string]*ComputeCapactiy{}

func init() {
	b, err := instanceTypesFile.ReadFile("aws_instances.json")
	if err != nil {
		panic(err)
		return
	}

	if err := json.Unmarshal(b, &instanceTypes); err != nil {
		panic(err)
	}
}

func Instances() []string {
	var arr []string

	for a, _ := range instanceTypes {
		arr = append(arr, a)
	}

	return arr
}

func InstanceCapacity(flavor string) *ComputeCapactiy {
	c, ok := instanceTypes[flavor]
	if !ok {
		return nil
	}

	return c
}

type ComputeCapactiy struct {
	CreditRate          float64
	MaxCredits          float64
	VCPUs               float64
	BaselineUtilization float64
}

type BurstCPU struct {
	Capacity    *ComputeCapactiy
	Utilisation float64
}

type BurstWindow struct {
	Name           string
	Utilisation    float64
	CreditBurnRate float64
	TimeToExhaust  time.Duration
	ShortWindow    time.Duration
	LongWindow     time.Duration
}

func (b *BurstWindow) timeToExhaust(c *BurstCPU, utilisation float64, window time.Duration) (time.Duration, error) {
	credits := c.Capacity.MaxCredits
	debitRate := b.creditsBurnRate(c, utilisation, window)

	if debitRate < c.Capacity.CreditRate {
		return time.Duration(math.Inf(1)), errors.New("The instance will never run out of credits since it is underutilised.")
	}

	return time.Duration((credits / debitRate) * dToFS(time.Hour)), nil
}

func (b *BurstWindow) creditsBurnRate(c *BurstCPU, utilisation float64, window time.Duration) float64 {
	creditsBurned := utilisation * float64(c.Capacity.VCPUs) * (dToFS(window) / dToFS(time.Minute))
	return creditsBurned / (dToFS(window) / dToFS(time.Hour))
}

func NewBurstCPU(cc *ComputeCapactiy, used float64) (*BurstCPU, error) {
	return &BurstCPU{
		Capacity:    cc,
		Utilisation: used,
	}, nil
}

type burstAlert struct {
}

func (b *burstAlert) String() string {
	return "burst alert"
}

const float64EqualityThreshold = 0.005

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= float64EqualityThreshold
}

func newBurstAlert(b *BurstCPU, name string, utilization float64) Alerter {
	carry := b.Capacity.MaxCredits - (24 * 60 * (b.Utilisation / 100.0) * b.Capacity.VCPUs)
	nMax := 24.0
	// (b.Capacity.MaxCredits - math.Max(carry, 0.0)) / b.Capacity.CreditRate
	nMin := math.Max(carry, 0.0) / (b.Capacity.VCPUs - (b.Capacity.CreditRate / 60))

	baseline := (b.Capacity.CreditRate * 100.0 / float64(b.Capacity.VCPUs)) / 60

	var cur float64
	for i := nMin; i <= nMax*60.0; i += 5 {
		u := ((carry + ((b.Capacity.CreditRate / 60) * i)) / (i * b.Capacity.VCPUs)) * 100
		if math.IsNaN(u) {
			continue
		}

		if almostEqual(cur, u) || u < 0 {
			break
		}

		if !almostEqual(u, 100) && !almostEqual(u, baseline*1.2) && !almostEqual(u, baseline*2) {
			continue
		}

		cur = u
		log.Println("utilization", u, "after minutes", i)
	}

	return &burstAlert{}
}

func BurstCalculator(b *BurstCPU) []Alerter {
	// Types of CPU burst alerts
	out := make([]Alerter, 2)

	// net-0 baseline per vcpu
	baseline := (b.Capacity.CreditRate * 100.0 / float64(b.Capacity.VCPUs)) / 60

	// A good starting point for a fast-burn threshold policy is 10x the
	// baseline with a short lookback period.
	// fastUtilization := math.Min(baseline*2, 100.0)

	// A good starting point for a slow-burn threshold is 2x the baseline with
	// a long lookback period.
	slowUtilization := math.Min(baseline*1.2, 100.0)

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
	out[0] = newBurstAlert(b, "slow", slowUtilization)

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
	// out[1] = newBurstAlert(b, "fast", fastUtilization)

	return out
}

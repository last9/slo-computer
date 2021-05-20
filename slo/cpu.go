package slo

import (
	"embed"
	"encoding/json"
	"fmt"
	"math"
	"time"
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

func NewBurstCPU(cc *ComputeCapactiy, used float64) (*BurstCPU, error) {
	return &BurstCPU{
		Capacity:    cc,
		Utilisation: used,
	}, nil
}

type burstAlert struct {
	utilization  float64
	exhaustAfter time.Duration
	alertAfter   time.Duration
}

func (b *burstAlert) String() string {
	return fmt.Sprintf(
		`
	Alert if %.2f %% consumption sustains for %s AND recent %s.
	At this rate, burst credits will deplete after %s
	`,
		b.utilization, b.alertAfter,
		maxD(b.alertAfter/12, 10*time.Minute), b.exhaustAfter,
	)
}

func newBurstAlert(b *BurstCPU, carry float64) []Alerter {
	// This is the equation that must be satisfied
	// - accured credits cannot exceed maxCredits
	// - When available credits will be less than the consumed credits based on
	// utilisation, it will breach.
	// min(maxCredits, carry + (creditRate * n)) <= (utilization/100) * vCPUs * n

	nMax := 24.0
	nMin := math.Max(carry, 0.0) /
		(b.Capacity.VCPUs - (b.Capacity.CreditRate / 60))

	// baseLine CPU utilization for this instance type at which rate of fill
	// = rate of depletion.
	baseline := (b.Capacity.CreditRate * 100.0 / float64(b.Capacity.VCPUs)) / 60

	var cur float64
	var alerts []Alerter

	// Since a CPU that would need more than 100, will only register as 100.
	// It's safe to assume that someone consuming ~ 100, can consume more aswell
	for i := nMin; i <= nMax*60.0; i += 5 {
		u := ((carry + ((b.Capacity.CreditRate / 60) * i)) /
			(i * b.Capacity.VCPUs)) *
			100

		if math.IsNaN(u) {
			continue
		}

		if almostEqual(cur, u) || u < 0 {
			break
		}

		ttl := time.Duration(int(i) * int(time.Minute))
		if almostEqual(u, 100) {
			alerts = append(alerts, &burstAlert{
				utilization:  u,
				exhaustAfter: ttl,
				alertAfter:   minD(ttl/2, time.Duration(1*time.Minute)),
			})
			continue
		}

		if !almostEqual(u, baseline*1.2) &&
			!almostEqual(u, baseline*2) {
			continue
		}

		alerts = append(alerts, &burstAlert{
			utilization:  u,
			exhaustAfter: ttl,
			alertAfter:   minD(ttl/2, time.Duration(60*time.Minute)),
		})
	}

	return alerts
}

func BurstCalculator(b *BurstCPU) []Alerter {
	// At the start of this minute, for the last 24 hours credits were accured.
	// but based on average consumption there were depletion too. Net sum is the
	// carry. This is only approximate. Since
	carry := b.Capacity.MaxCredits -
		(24 * 60 * (b.Utilisation / 100.0) * b.Capacity.VCPUs)

	return newBurstAlert(b, carry)
}

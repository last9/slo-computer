package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/last9/slo-computer/slo"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

const errorMessage = `
	If this service reported %.6f errors for a duration of %s
	SLO (for the entire duration) will be defeated within %s

	Probably
	- Use ONLY spike alert model, and not SLOs (easiest)
	- Reduce the MTTR for this service (toughest)
	- SLO is too aggressive and can be lowered (business decision)
	- Combine multiple services into one single service (team wide)
`

type suggestCmd struct {
	throughput float64
	sloDesire  float64
	sloPeriod  int
}

func (c *suggestCmd) run(ctx *kingpin.ParseContext) error {
	s, err := slo.NewSLO(
		time.Duration(c.sloPeriod)*time.Hour,
		c.throughput, c.sloDesire,
	)

	if err != nil {
		return err
	}

	imp, yes := slo.IsLowTraffic(s)
	if yes {
		return errors.Errorf(
			errorMessage, imp.Errors, imp.Duration,
			imp.BreaksAfter,
		)
	}

	a := slo.AlertCalculator(s)
	for _, aw := range a {
		fmt.Println(aw)
	}

	return nil
}

func suggestCommand(app *kingpin.Application) {
	c := &suggestCmd{}
	sg := app.Command("suggest", "suggest alerts based on the input").Action(c.run)

	sg.Flag("throughput", "Throughput for this service").Required().FloatVar(&c.throughput)
	sg.Flag("slo", "Desired SLO for this service").Required().FloatVar(&c.sloDesire)
	sg.Flag("duration", "Duration for the SLO").Required().IntVar(&c.sloPeriod)
}

type burstCPUCmd struct {
	instanceType   string
	cpuUtilization float64
	timeDuration      int
}

func (c *burstCPUCmd) run(ctx *kingpin.ParseContext) error {
	//TODO: Validate Input & compute permisisble burst
	b, err := slo.NewBurstCPU(
		c.instanceType,
		c.cpuUtilization,
		time.Duration(c.timeDuration)*time.Hour,
	)
	if err != nil {
		return err
	}

	log.Println(b)

	return nil
}

func burstCPUCommand(app *kingpin.Application) {
	c := &burstCPUCmd{}
	sg := app.Command("cpu-burst", "compute permissible burst interval for an instance").Action(c.run)

	sg.Flag("instance", "Instance type").Required().StringVar(&c.instanceType)
	sg.Flag("utilisation", "Average CPU Utilisation").Required().FloatVar(&c.cpuUtilization)
	sg.Flag("duration", "Duration (in hours) for the utilization").Required().IntVar(&c.timeDuration)
}

func main() {
	app := kingpin.New("slo", "Last9 SLO toolkit")
	suggestCommand(app)
	burstCPUCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

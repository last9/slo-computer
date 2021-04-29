package main

import (
	"fmt"
	"os"
	"time"

	"github.com/last9/slo-computer/slo"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

const errorMessage = `
	If this service reported %.3f errors for a duration of %s
	SLO (for the entire duration) will be defeated wihin %s

	Probably
	- Use ONLY spike alert model, and not SLOs (easiest)
	- Reduce the MTTR for this service (toughest)
	- SLO is too aggressive and can be lowerd (business decision)
	- Combine multiple services into one single service (teamwide)
`

type suggestCmd struct {
	throughput float64
	slo_desire float64
	slo_period int
}

func (c *suggestCmd) run(ctx *kingpin.ParseContext) error {
	s, err := slo.NewSLO(
		time.Duration(time.Duration(c.slo_period)*time.Hour),
		c.throughput, c.slo_desire,
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
	sg.Flag("slo", "Desired SLO for this service").Required().FloatVar(&c.slo_desire)
	sg.Flag("duration", "Duration for the SLO").Required().IntVar(&c.slo_period)
}

func main() {
	app := kingpin.New("slo", "Last9 SLO toolkit")
	suggestCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))

}

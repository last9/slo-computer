package main

import (
	"fmt"
	"time"

	"github.com/last9/slo-computer/slo"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

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
		fmt.Println(aw.String())
	}

	return nil
}

func suggestCommand(app *kingpin.Application) {
	c := &suggestCmd{}
	sg := app.Command(
		"suggest",
		"suggest alerts based on service throughput and SLO duration",
	).Action(c.run)

	sg.Flag("throughput", "Throughput for this service").Required().FloatVar(&c.throughput)
	sg.Flag("slo", "Desired SLO for this service").Required().FloatVar(&c.sloDesire)
	sg.Flag("duration", "Duration for the SLO").Required().IntVar(&c.sloPeriod)
}

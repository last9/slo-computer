package main

import (
	"fmt"

	"github.com/last9/slo-computer/slo"
	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

type burstCPUCmd struct {
	instanceType   string
	cpuUtilization float64
}

func (c *burstCPUCmd) suggest(ctx *kingpin.ParseContext) error {
	cc := slo.InstanceCapacity(c.instanceType)
	if cc == nil {
		return errors.Errorf("unrecognized instance: %v", c.instanceType)
	}

	if c.cpuUtilization > 100 || c.cpuUtilization < 0 {
		return errors.Errorf("avg cpu should be: 0 <= %v <= 100", c.cpuUtilization)
	}

	//TODO: Validate Input & compute permisisble burst
	b, err := slo.NewBurstCPU(cc, c.cpuUtilization)
	if err != nil {
		return err
	}

	a := slo.BurstCalculator(b)
	for _, aw := range a {
		if aw == nil {
			continue
		}
		fmt.Println(aw.String())
	}
	return nil
}

func burstCPUCommand(app *kingpin.Application) {
	c := &burstCPUCmd{}
	sg := app.Command(
		"cpu-suggest", "compute permissible burst interval for an instance",
	).Action(c.suggest)

	sg.Flag("instance", "Instance type").Required().EnumVar(
		&c.instanceType, slo.Instances()...,
	)

	sg.Flag(
		"utilization",
		"Average utilization over last 24 hours",
	).Required().Float64Var(&c.cpuUtilization)
}

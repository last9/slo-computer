package main

import (
	"fmt"
	"log"
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

var (
	throughput = kingpin.Flag("throughput", "Throughput for this service").Required().Float()
	slo_desire = kingpin.Flag("slo", "Desired SLO for this service").Required().Float()
	slo_period = kingpin.Flag("duration", "Duration for the SLO").Required().Int()
)

func main() {
	kingpin.Parse()

	s, err := slo.NewSLO(
		time.Duration(time.Duration(*slo_period)*time.Hour),
		*throughput, *slo_desire,
	)

	if err != nil {
		log.Fatal(err)
	}

	imp, yes := slo.IsLowTraffic(s)
	if yes {
		log.Fatal(errors.Errorf(
			errorMessage, imp.Errors, imp.Duration,
			imp.BreaksAfter,
		))
	}

	a := slo.AlertCalculator(s)
	for _, aw := range a {
		fmt.Println(aw)
	}
}

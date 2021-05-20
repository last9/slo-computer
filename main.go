package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var Version = "0.0.2"

const errorMessage = `
	If this service reported %.6f errors for a duration of %s
	SLO (for the entire duration) will be defeated within %s

	Probably
	- Use ONLY spike alert model, and not SLOs (easiest)
	- Reduce the MTTR for this service (toughest)
	- SLO is too aggressive and can be lowered (business decision)
	- Combine multiple services into one single service (team wide)
`

func main() {
	app := kingpin.New("slo", "Last9 SLO toolkit")
	app = app.Version(Version)
	suggestCommand(app)
	burstCPUCommand(app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

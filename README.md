<a href="https://last9.io"><img src="https://last9.github.io/assets/last9-github-badge.svg" align="right" /></a>

# SLO Computer

> [!Note]
> @last9 advocates using Service Level Objectives.
> One of the biggest challenges we run into is the lack of practical algorithms behind Burn Rate and alerting. This is our first attempt at it.

## What is SLO Computer?

SLO Computer simplifies the complex world of Service Level Objectives (SLOs), error budgets, and alerting. 

SLOs, error windows, burn rates, and budget spend are convoluted terms that can throw anyone off. Even the SRE workbook by Google can leave you with a lot of open questions. We continue to be amazed by how widely misunderstood this topic is (and how easy it can make your lives if used well).

This toolkit helps SREs and DevOps engineers:

- Calculate appropriate alert thresholds based on service throughput and desired SLO targets
- Determine if a service has enough traffic to benefit from SLO-based alerting
- Generate alert policies for AWS burstable CPU instances

## Installation and Setup

### Prerequisites
- Go 1.16 or later

### Installation Options

#### Using Go
```bash
# Install directly using Go
go install github.com/last9/slo-computer@latest
```

#### Using Docker
```bash
# Pull the Docker image
docker pull last9/slo-computer:latest

# Run using Docker
docker run last9/slo-computer:latest --help
```

#### Building from Source
```bash
# Clone the repository
git clone https://github.com/last9/slo-computer.git
cd slo-computer

# Build using Make
make build
```

### Quick Start

The project includes a Makefile with helpful commands:

```bash
# Build the application
make build

# Run tests
make test

# Run an example service SLO calculation
make example-service

# Run an example CPU burst calculation
make example-cpu

# See all available commands
make help
```

## Usage

```bash
usage: slo [<flags>] <command> [<args> ...]

Last9 SLO toolkit

Flags:
  --help                Show context-sensitive help (also try --help-long and --help-man).
  --version             Show application version.
  --config=CONFIG       Path to configuration file
  --output=FORMAT       Output format (text, json, yaml)

Commands:
  help [<command>...]
    Show help.

  suggest --throughput=THROUGHPUT --slo=SLO --duration=DURATION
    suggest alerts based on service throughput and SLO duration

  cpu-suggest --instance=INSTANCE --utilization=UTILIZATION
    suggest alerts based on CPU utilization and Instance type
```

### Command Parameters

#### `suggest` Command
- `--throughput`: Number of requests per minute your service handles
- `--slo`: Your desired SLO percentage (e.g., 99.9)
- `--duration`: SLO time period in hours (e.g., 720 for 30 days)

#### `cpu-suggest` Command
- `--instance`: AWS instance type (e.g., t3.micro, t3a.xlarge)
- `--utilization`: Average CPU utilization percentage (0-100)

### Using Configuration Files

You can define your services and configurations in YAML or JSON files:

```yaml
# slo-config.yaml
services:
  api-gateway:
    throughput: 4200
    slo: 99.9
    duration: 720
  
  background-processor:
    throughput: 100
    slo: 99.5
    duration: 168

cpus:
  web-server:
    instance: t3a.xlarge
    utilization: 15
```

Then use it with:

```bash
# For a specific service
./slo-computer suggest --config=slo-config.yaml --service=api-gateway

# For a specific CPU
./slo-computer cpu-suggest --config=slo-config.yaml --service=web-server
```

### Output Formats

SLO Computer supports multiple output formats:

```bash
# Default text output
./slo-computer suggest --throughput=4200 --slo=99.9 --duration=720

# JSON output
./slo-computer suggest --throughput=4200 --slo=99.9 --duration=720 --output=json

# YAML output
./slo-computer suggest --throughput=4200 --slo=99.9 --duration=720 --output=yaml
```

## CI/CD Integration

### GitHub Actions

You can use SLO Computer in your GitHub Actions workflows:

```yaml
name: SLO Analysis

on:
  schedule:
    - cron: '0 0 * * 1'  # Weekly on Monday
  workflow_dispatch:  # Manual trigger

jobs:
  analyze-slos:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Run SLO Computer
        uses: last9/slo-computer-action@v1
        with:
          command: suggest
          config-file: .github/slo-config.yaml
          service-name: api-gateway
          output-format: json
```

### Docker Integration

```bash
# Mount your config file and run
docker run -v $(pwd)/slo-config.yaml:/config.yaml last9/slo-computer:latest suggest --config=/config.yaml --service=api-gateway
```

## Examples

### Service SLO Alerts

**Q: What alerts should I set for my service to achieve 99.9% availability over 30 days?**

```bash
./slo-computer suggest --throughput=4200 --slo=99.9 --duration=720
```

Output:
```
Alert if error_rate > 0.002 for last [24h0m0s] and also last [2h0m0s]
This alert will trigger once 6.67% of error budget is consumed,
and leaves 360h0m0s before the SLO is defeated.


Alert if error_rate > 0.010 for last [1h0m0s] and also last [5m0s]
This alert will trigger once 1.39% of error budget is consumed,
and leaves 72h0m0s before the SLO is defeated.
```

JSON Output:
```json
[
  {
    "type": "slow_burn",
    "error_rate": 0.002,
    "long_window": "24h0m0s",
    "short_window": "2h0m0s",
    "budget_consumed": 0.0667,
    "time_remaining": "360h0m0s"
  },
  {
    "type": "fast_burn",
    "error_rate": 0.01,
    "long_window": "1h0m0s",
    "short_window": "5m0s",
    "budget_consumed": 0.0139,
    "time_remaining": "72h0m0s"
  }
]
```

**Q: What about a low-traffic service?**

```bash
./slo-computer suggest --throughput=100 --slo=99.9 --duration=168
```

Output:
```
slo-computer: error:
	If this service reported 10.000 errors for a duration of 5m0s
	SLO (for the entire duration) will be defeated wihin 1h40m47s

	Probably
	- Use ONLY spike alert model, and not SLOs (easiest)
	- Reduce the MTTR for this service (toughest)
	- SLO is too aggressive and can be lowered (business decision)
	- Combine multiple services into one single service (team wide)
```

### CPU Burst Credit Alerts

**Q: What alerts should I set for my AWS burstable instance?**

```bash
./slo-computer cpu-suggest --instance=t3a.xlarge --utilization=15
```

Output:
```
Alert if 100.00 % consumption sustains for 10m0s AND recent 5m0s.
At this rate, burst credits will deplete after 10h0m0s


Alert if 80.00 % consumption sustains for 3h45m0s AND recent 55m0s.
At this rate, burst credits will deplete after 15h0m0s
```

## Understanding the Results

### For Service SLOs

The tool generates two types of alerts:
1. **Slow burn alert**: Detects gradual error rate increases that would eventually exhaust your error budget
2. **Fast burn alert**: Detects sudden spikes in error rates that require immediate attention

Each alert includes:
- The error rate threshold to monitor
- The time windows to evaluate
- How much of your error budget would be consumed when the alert triggers
- How much time remains before your SLO is breached if the error rate continues

### For CPU Burst Credits

The tool generates alerts that help you monitor when your AWS burstable instance might run out of CPU credits:
- Alert thresholds for different CPU utilization levels
- Time windows to monitor
- Time until credit depletion at the current rate

## Key Concepts

### Service SLOs
- **Throughput**: The number of requests your service handles per minute
- **SLO**: Your Service Level Objective (e.g., 99.9% availability)
- **Duration**: The time period for your SLO in hours (e.g., 720 for 30 days)
- **Error Budget**: The amount of allowable errors within your SLO period (calculated as `(100% - SLO%) * total requests`)
- **Burn Rate**: How quickly you're consuming your error budget relative to the expected rate

### CPU Burst Credits
- **Instance**: AWS burstable instance type (T2, T3, T4g families)
- **Utilization**: Average CPU utilization percentage
- **Credit Rate**: How quickly the instance earns CPU credits
- **Baseline Performance**: The CPU performance level the instance can sustain indefinitely

## Using as a Library

You can also use SLO Computer as a library in your Go projects:

```go
import (
    "time"
    "github.com/last9/slo-computer/slo"
)

// Create a new SLO
s, err := slo.NewSLO(
    time.Duration(720)*time.Hour, // SLO period of 30 days
    4200,                         // 4200 requests per minute
    99.9,                         // 99.9% availability target
)

// Calculate alerts
alerts := slo.AlertCalculator(s)

// For CPU burst calculations
cc := slo.InstanceCapacity("t3.micro")
b, err := slo.NewBurstCPU(cc, 75.0) // 75% utilization
burstAlerts := slo.BurstCalculator(b)
```

## Troubleshooting

### Common Errors

**Error: "strconv.ParseFloat: parsing "SLO": invalid syntax"**  
Make sure to replace "SLO" with an actual number (e.g., 99.9) in your command:
```bash
# Incorrect
./slo-computer suggest --throughput=1000000 --slo=SLO --duration=90

# Correct
./slo-computer suggest --throughput=1000000 --slo=99.9 --duration=90
```

**Error about low traffic services**  
If you receive a message about your service being low-traffic, consider:
- Using spike-based alerting instead of SLO-based alerting
- Combining multiple services to increase the traffic volume
- Lowering your SLO target to a more achievable level

## Roadmap

We're actively working on improving SLO Computer. Check out our roadmap:
- [Open Issues](OPEN_ISSUES.md) - Planned improvements and bug fixes
- [Feature Enhancements](FEATURES.md) - Upcoming features and user experience improvements

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

# About Last9

This project is sponsored and maintained by [Last9](https://last9.io). Last9 is a telemetry data platform.

<a href="https://last9.io"><img src="https://last9.github.io/assets/email-logo-green.png" alt="" loading="lazy" height="40px" /></a>

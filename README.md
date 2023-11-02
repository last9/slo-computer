# slo-computer

> [!Note]
> Around two years ago @last9 started advocating using Service Level Objectives.
> One of the biggest challenges was the lack of practical algorithms behind Burn Rate and alerting. This was our first attempt at it.
> If you would like us to release these algorithms, go ahead and help us reach 250 stars ⭐️. 

SLOs, Error windows and alerts are complicated. Here's an attempt to make it
easy.

SLO, burn_rate, error_rate, budget_spend are convoluted terms that can throw one
off. Even the SRE workbook by Google can leave you with a lot of open questions.

The concept of SLOs and SLIs has existed for a long time now, but we continue to
be amazed by how widely misunderstood this topic is (and how easy it can make
your lives if used well).

We are building a sandbox for our DevOps and SRE community - SLO computer - a
product that makes setting and monitoring SLOs for all your services intuitively
seamless and blazingly fast.

### Usage

```bash
usage: slo [<flags>] <command> [<args> ...]

Last9 SLO toolkit

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.

  suggest --throughput=THROUGHPUT --slo=SLO --duration=DURATION
    suggest alerts based on service throughput and SLO duration

  cpu-suggest --instance=INSTANCE --utilization=UTILIZATION
    suggest alerts based on CPU utilization and Instance type
```

The goal of this command (has an importable lib too) is to factor in some "bare
minimum" input to

- Is this a Low traffic service in which case it makes little sense to use an
  SLO approach
- Compute the _actual_ alert values and condition to set alerts on

### Examples

**Q: What alerts should I set for my service to achieve 99.9 % availability over
30 days**

```bash
✗ ./slo-computer suggest --throughput=4200 --slo=99.9 --duration=720

		Alert if error_rate > 0.002 for last [24h0m0s] and also last [2h0m0s]
		This alert will trigger once 6.67% of error budget is consumed,
		and leaves 360h0m0s before the SLO is defeated.


		Alert if error_rate > 0.010 for last [1h0m0s] and also last [5m0s]
		This alert will trigger once 1.39% of error budget is consumed,
		and leaves 72h0m0s before the SLO is defeated.
```

**Q: What alerts should I set for my service with throughpput 100rpm to achieve
90 % availability over 7 days**

```bash
✗ ./slo-computer suggest --throughput=100 --slo=99.9 --duration=168
slo-computer: error:
	If this service reported 10.000 errors for a duration of 5m0s
	SLO (for the entire duration) will be defeated wihin 1h40m47s

	Probably
	- Use ONLY spike alert model, and not SLOs (easiest)
	- Reduce the MTTR for this service (toughest)
	- SLO is too aggressive and can be lowerd (business decision)
	- Combine multiple services into one single service (teamwide)
, try --help
```

**Q: What alerts should I set for my burst CPU**

```bash
✗ ./slo-computer cpu-suggest --instance=t3a.xlarge --utilization=15

	Alert if 100.00 % consumption sustains for 10m0s AND recent 5m0s.
	At this rate, burst credits will deplete after 10h0m0s


	Alert if 80.00 % consumption sustains for 3h45m0s AND recent 55m0s.
	At this rate, burst credits will deplete after 15h0m0s
```

---

# About Last9

This project is sponsored and maintained by [Last9](https://last9.io). Last9 builds reliability tools for SRE and DevOps.

<a href="https://last9.io"><img src="https://last9.github.io/assets/email-logo-green.png" alt="" loading="lazy" height="40px" /></a>

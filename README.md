# slo-computer
SLOs, Error windows and alerts are complicated. Here an attempt to make it easy

SLO, burn_rate, error_rate, budget_spend are convoluted terms that can throw one off.
Even the SRE workbook by Google can leave you with a lot of open questions.

# Goal

The goal of this command (has an importable lib too) is to factor in some "bare minimum" input to
- Is this a Low traffic service in which case it makes little sense to use an SLO approach
- Compute the *actual* alert values and condition to set alerts on

## Examples

**Q: What alerts should I set for my service to achieve 99.9 % availability over 30 days**
```bash 
➜  slo-computer git:(master) ✗ ./slo-computer suggest --throughput=4200 --slo=99.9 --duration=720

		Alert if error_rate > 0.002 for last [24h0m0s] and also last [2h0m0s]
		This alert will trigger once 6.67% of error budget is consumed,
		and leaves 360h0m0s before the SLO is defeated.
		

		Alert if error_rate > 0.010 for last [1h0m0s] and also last [5m0s]
		This alert will trigger once 1.39% of error budget is consumed,
		and leaves 72h0m0s before the SLO is defeated.
```

**Q: What alerts should I set for my service with throughpput 100rpm to achieve 90 % availability over 7 days**

```bash
➜  slo-computer git:(master) ✗ ./slo-computer suggest --throughput=100 --slo=99.9 --duration=168
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

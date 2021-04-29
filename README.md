# slo-computer
SLOs, Error windows and alerts are complicated. Here an attempt to make it easy


`➜  slo-computer git:(master) ./main --throughput=4200 --slo=99.9 --duration=720`

```bash
		Alert if error_rate > 0.002 for last [24h0m0s] and also last [2h0m0s]
		This alert will trigger once 6.67% of error budget is consumed,
		and leaves 360h0m0s before the SLO is defeated.
		

		Alert if error_rate > 0.010 for last [1h0m0s] and also last [5m0s]
		This alert will trigger once 1.39% of error budget is consumed,
		and leaves 72h0m0s before the SLO is defeated.
```

`➜  slo-computer git:(master) ./main --throughput=4200 --slo=99.9 --duration=5`

```bash
		Alert if error_rate > 0.002 for last [2h0m0s] and also last [10m0s]
		This alert will trigger once 80.00% of error budget is consumed,
		and leaves 2h30m0s before the SLO is defeated.
		

		Alert if error_rate > 0.010 for last [5m0s] and also last [2m0s]
		This alert will trigger once 16.67% of error budget is consumed,
		and leaves 30m0s before the SLO is defeated.
```

# Metrics

The signare exposes key metrics to measure its state through a `Prometheus` server. This document describes the complete set of metrics provided by the application.

The target audience of this document are signare operators and developers.

The following metric types are used depending on the value that needs to be measured:

- Counter: Cumulative metric that represents a single monotonically increasing counter.

For more details on types of metrics, please refer to [Prometheus metrics documentation](https://prometheus.io/docs/concepts/metric_types/>)

## HTTP metrics

| Name                        | Labels             | Type    | Description                                                     |
|-----------------------------|--------------------|---------|-----------------------------------------------------------------|
| **forbidden_access_count**  | "action", "error"  | counter | Total number of attempts to perform a given unauthorized action |


## Default GO process metrics exposed by Prometheus GO client library

The default set of metrics that Prometheus’ [client_golang](https://github.com/prometheus/client_golang>) exposes are also exposed by the signare.

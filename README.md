# Qubership Prometheus Adapter Operator

This repository contains the Kubernetes Operator for
[Prometheus Adapter](https://github.com/kubernetes-sigs/prometheus-adapter).

Prometheus Adapter for Kubernetes Metrics API is implementation of the Kubernetes Custom, Resource and External
[Metric APIs](https://github.com/kubernetes/metrics).

This adapter is therefore suitable for use with the autoscaling/v2 `Horizontal Pod Autoscaler` in Kubernetes 1.6+.
It can also replace the [metrics server](https://github.com/kubernetes-incubator/metrics-server) on clusters
that already run Prometheus and collect the appropriate metrics.

## Support Kubernetes Version

The `prometheus-operator-adapter` and `prometheus-adapter` can be installed in `Kubernetes 1.16+`.

## Documentation

* [Installation](docs/install.md)
* CustomResource API descriptions:
  * [PrometheusAdapter](docs/api.md)
  * [CustomScaleMetricRule](docs/api.md)

## Before you begin

To deploy the `qubership-prometheus-adapter-operator` you need to have:

* kubectl 1.16+
* helm 3+
* configured `kubectl` context on your cloud

To build `qubership-prometheus-adapter-operator` you need to have:

* `make`
* `gcc`
* `golang` 1.19+
* `operator-sdk` 1.x

## Installation

This repository contains the Helm chart which allow to deploy operator in Kubernetes.
Also the copy of this Helm chart include in monitoring-operator Helm chart. And the prometheus-adapter can be
deployed with monitoring.

To install the qubership-prometheus-adapter-operator without monitoring (and for example connect to already deployed Prometheus),
execute the following steps:

1. Checkout this repository or download archive with repository content

2. Fill values for deploy (either create a new `custom-value.yaml` or change default `value.yaml`).
   All available values described in [Installation](docs/install.md) and in
   [values.yaml](charts/qubership-prometheus-adapter-operator/values.yaml).
3. Run deploy with using Helm

   ```bash
   helm install <release_name> --namespace=<namespace> charts/qubership-prometheus-adapter-operator -f /path/to/custom-value.yaml
   ```

## Release Images

All release images for the `qubership-prometheus-adapter-operator` and for the `prometheus-adapter` published in

* [ghcr.io/netcracker/qubership-prometheus-adapter](https://github.com/Netcracker/qubership-prometheus-adapter-operator/pkgs/container/qubership-prometheus-adapter)
* [ghcr.io/netcracker/qubership-prometheus-adapter-operator](https://github.com/Netcracker/qubership-prometheus-adapter-operator/pkgs/container/qubership-prometheus-adapter-operator)

## Configuration

## Example

For try to use horizontal pod autoscaling, see necessary configs

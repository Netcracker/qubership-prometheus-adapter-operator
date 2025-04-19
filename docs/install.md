# Installation guide

## TL;DR

```console
helm install prometheus-adapter charts/qubership-prometheus-adapter-operator
```

## Prerequisites

- Kubernetes v1.15+

## Installing the Chart

To install the chart with the release name `prometheus-adapter`:

```console
helm install prometheus-adapter charts/qubership-prometheus-adapter-operator
```

The command deploys a Stash operator on the Kubernetes cluster in the default configuration.
The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `prometheus-adapter`:

```console
helm uninstall prometheus-adapter
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the `prometheus-adapter` chart and their default values.

<!-- markdownlint-disable line-length -->

| Parameter                                        | Description                                                                                                                                                                                                                                                                                                                                                                                             | Default                                                             |
| ------------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------- |
| nameOverride                                     | Provide a name in place of qubership-prometheus-adapter-operator for `app:` labels                                                                                                                                                                                                                                                                                                                      | `""`                                                                |
| fullnameOverride                                 | Provide a name to substitute for the full names of resources                                                                                                                                                                                                                                                                                                                                            | `""`                                                                |
| image                                            | A docker image to use for qubership-prometheus-adapter-operator deployment Type: string Mandatory: yes                                                                                                                                                                                                                                                                                                  | `"ghcr.io/netcracker/qubership-prometheus-adapter-operator:latest"` |
| role.install                                     | Allow to disable create Role during deploy Type: object Mandatory: no                                                                                                                                                                                                                                                                                                                                   | `yes`                                                               |
| clusterRole.install                              | Allow to disable create ClusterRole during deploy Type: object Mandatory: no                                                                                                                                                                                                                                                                                                                            | `yes`                                                               |
| roleBinding.install                              | Allow to disable create RoleBinding during deploy Type: object Mandatory: no                                                                                                                                                                                                                                                                                                                            | `yes`                                                               |
| clusterRoleBinding.install                       | Allow to disable create ClusterRoleBinding during deploy Type: object Mandatory: no                                                                                                                                                                                                                                                                                                                     | `yes`                                                               |
| securityContext                                  | SecurityContext holds pod-level security attributes. The parameters are required if a Pod Security Policy is enabled for Kubernetes cluster and required if a Security Context Constraints is enabled for Openshift cluster. Mandatory: no                                                                                                                                                              | `{}`                                                                |
| resources                                        | A special supplemental group that applies to all containers in a pod. If unset, the Kubelet will not modify the ownership and permissions of any volume. Mandatory: no fsGroup: 2000 The resources describes the compute resource requests and limits for single Pods. Ref: <https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#example-1> Type: object Mandatory: no       | `{}`                                                                |
| nodeSelector                                     | Allow define which Nodes the Pods are scheduled on. Type: map[string] Mandatory: no Default: not set                                                                                                                                                                                                                                                                                                    | `{}`                                                                |
| APIService.resourceMetrics                       | Enable/disable creating APIServices for `metrics.k8s.io`. Type: bool Mandatory: no Default: true                                                                                                                                                                                                                                                                                                        | `true`                                                              |
| APIService.customMetrics                         | Enable/disable creating APIServices for `custom.metrics.io`. Type: bool Mandatory: no Default: true                                                                                                                                                                                                                                                                                                     | `true`                                                              |
| prometheusAdapter.image                          | A docker image to use for prometheus-adapter deployment Type: string Mandatory: yes                                                                                                                                                                                                                                                                                                                     | `"ghcr.io/netcracker/prometheus-adapter:latest"`                    |
| prometheusAdapter.metricsRelistInterval          | This is the interval at which to update the cache of available metrics from Prometheus. Since the adapter only lists metrics during discovery that exist between the current time and the last discovery query, your relist interval should be equal to or larger than your Prometheus scrape interval, otherwise your metrics will occasionally disappear from the adapter. Type: string Mandatory: no | `"1m"`                                                              |
| prometheusAdapter.enableResourceMetrics          | Enable adapter for `metrics.k8s.io`. Type: bool                                                                                                                                                                                                                                                                                                                                                         | `false`                                                             |
| prometheusAdapter.enableCustomMetrics            | Enable adapter for `custom.metrics.k8s.io`. Type: bool                                                                                                                                                                                                                                                                                                                                                  | `true`                                                              |
| prometheusAdapter.prometheusUrl                  | This is the URL used to connect to Prometheus. It will eventually contain query parameters to configure the connection. Type: string Mandatory: no                                                                                                                                                                                                                                                      | `""`                                                                |
| prometheusAdapter.customScaleMetricRulesSelector | CustomResources's labels to match for CustomScaleMetricRules discovery. If nil, only check all namespaces. Type: LabelSelector Mandatory: no                                                                                                                                                                                                                                                            | `[]`                                                                |
| prometheusAdapter.securityContext                | SecurityContext holds pod-level security attributes. The parameters are required if a Pod Security Policy is enabled for Kubernetes cluster and required if a Security Context Constraints is enabled for Openshift cluster. Mandatory: no                                                                                                                                                              | `{}`                                                                |
| prometheusAdapter.resources                      | A special supplemental group that applies to all containers in a pod. If unset, the Kubelet will not modify the ownership and permissions of any volume. Mandatory: no fsGroup: 2000 The resources describes the compute resource requests and limits for single Pods. Ref: <https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#example-1> Type: object Mandatory: no       | `{}`                                                                |
| prometheusAdapter.nodeSelector                   | Allow define which Nodes the Pods are scheduled on. Type: map[string] Mandatory: no Default: not set                                                                                                                                                                                                                                                                                                    | `{}`                                                                |

<!-- markdownlint-enable line-length -->

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
helm install prometheus-adapter charts/qubership-prometheus-adapter-operator --set image="ghcr.io/netcracker/qubership-prometheus-adapter-operator:latest"
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
helm install prometheus-adapter charts/qubership-prometheus-adapter-operator --values values.yaml
```

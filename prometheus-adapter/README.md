# Prometheus Adapter

This repository builds a custom Docker image of `prometheus-adapter` to add an ability to reload config without the pod's
restart.

The `entrypoint.sh` script logic:

* wrap the run of `prometheus-adapter` binary
* watch the directory with the config
* re-run `prometheus-adapter` binary if the configuration file was changed

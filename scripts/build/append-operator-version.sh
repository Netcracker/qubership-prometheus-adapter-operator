#!/usr/bin/env bash

if [[ "$OSTYPE" == "darwin"* ]]; then
    # Add the custom annotation
    find charts/qubership-prometheus-adapter-operator/crds -name '*.yaml' -exec sed -i '' -e "/^    controller-gen.kubebuilder.io.version.*/a\\
        qubership-prometheus-adapter-operator.monitoring.qubership.org/version: $VERSION" {} +
    
    # Add the labels right after 'annotations:'
    find charts/qubership-prometheus-adapter-operator/crds -name '*.yaml' -exec sed -i '' -e "/^  annotations:/a\\
      labels:\\
          app.kubernetes.io/component: qubership-prometheus-adapter-operator\\
          app.kubernetes.io/part-of: monitoring" {} +
else
    # Linux
    find charts/qubership-prometheus-adapter-operator/crds -name '*.yaml' -exec sed -i "/^    controller-gen.kubebuilder.io.version.*/a\\
        qubership-prometheus-adapter-operator.monitoring.qubership.org/version: $VERSION" {} +
    
    find charts/qubership-prometheus-adapter-operator/crds -name '*.yaml' -exec sed -i "/^  annotations:/a\\
      labels:\\
          app.kubernetes.io/component: qubership-prometheus-adapter-operator\\
          app.kubernetes.io/part-of: monitoring" {} +
fi

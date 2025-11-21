// Copyright 2025 NetCracker Technology Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prometheusadapter

import (
	promv1 "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	"github.com/go-logr/logr"
	pmodel "github.com/prometheus/common/model"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// PrometheusAdapterManager holds methods to control prometheus-adapter settings, e.g. ConfigMap changing.
// We want this kind of entity because we need to change prometheus-adapter settings in several controllers.
type PrometheusAdapterManager struct {
	client client.Client
	log    logr.Logger
}

// NewPrometheusAdapterManager creates an instance of PrometheusAdapterManager.
func NewPrometheusAdapterManager(client client.Client, log logr.Logger) *PrometheusAdapterManager {
	return &PrometheusAdapterManager{
		client: client,
		log:    log.WithName("manager-prometheusadapter"),
	}
}

// ConfigMap represents prometheus-adapter ConfigMap which holds metrics rules.
type ConfigMap struct {
	CustomMetricRules []promv1.CustomMetricRuleConfig `json:"rules,omitempty"`
	ResourceRules     ResourceRules                   `json:"resourceRules,omitempty"`
}

// ResourceRules describe the rules for querying resource metrics
// API results.  It's assumed that the same metrics can be used
// to aggregate across different resources.
type ResourceRules struct {
	CPU    ResourceRule `json:"cpu" yaml:"cpu"`
	Memory ResourceRule `json:"memory" yaml:"memory"`
	// Window is the window size reported by the resource metrics API.  It should match the value used
	// in your containerQuery and nodeQuery if you use a `rate` function.
	Window pmodel.Duration `json:"window" yaml:"window"`
}

// ResourceRule describes how to query metrics for some particular
// system resource metric.
type ResourceRule struct {
	// Container is the query used to fetch the metrics for containers.
	ContainerQuery string `json:"containerQuery" yaml:"containerQuery"`
	// NodeQuery is the query used to fetch the metrics for nodes
	// (for instance, simply aggregating by node label is insufficient for
	// cadvisor metrics -- you need to select the `/` container).
	NodeQuery string `json:"nodeQuery" yaml:"nodeQuery"`
	// Resources specifies how associated Kubernetes resources should be discovered for
	// the given metrics.
	Resources ResourceMapping `json:"resources" yaml:"resources"`
	// ContainerLabel indicates the name of the Prometheus label containing the container name
	// (since "container" is not a resource, this can't go in the `resources` block, but is similar).
	ContainerLabel string `json:"containerLabel" yaml:"containerLabel"`
}

// ResourceMapping specifies how to map Kubernetes resources to Prometheus labels
type ResourceMapping struct {
	// Template specifies a golang string template for converting a Kubernetes
	// group-resource to a Prometheus label.  The template object contains
	// the `.Group` and `.Resource` fields.  The `.Group` field will have
	// dots replaced with underscores, and the `.Resource` field will be
	// singularized.  The delimiters are `<<` and `>>`.
	Template string `json:"template,omitempty" yaml:"template,omitempty"`
	// Overrides specifies exceptions to the above template, mapping label names
	// to group-resources
	Overrides map[string]GroupResource `json:"overrides,omitempty" yaml:"overrides,omitempty"`
	// Namespaced ignores the source namespace of the requester and requires one in the query
	Namespaced *bool `json:"namespaced,omitempty" yaml:"namespaced,omitempty"`
}

// GroupResource represents a Kubernetes group-resource.
type GroupResource struct {
	Group    string `json:"group,omitempty" yaml:"group,omitempty"`
	Resource string `json:"resource" yaml:"resource"`
}

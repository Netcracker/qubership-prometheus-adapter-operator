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

package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CustomScaleMetricRuleSpec defines the desired state of CustomScaleMetricRule
type CustomScaleMetricRuleSpec struct {
	Rules []CustomMetricRuleConfig `json:"rules,omitempty"`
}

// CustomMetricRuleConfig defines the metric exposing rule from Prometheus.
// This structure is similar to the DiscoveryRule from github.com/directxman12/k8s-prometheus-adapter/pkg/config
// but we can not use the original structure because it is not compliant with kube-builder's CRD generator.
type CustomMetricRuleConfig struct {
	// SeriesQuery specifies which metrics this rule should consider via a Prometheus query
	// series selector query.
	SeriesQuery string `json:"seriesQuery"`
	// SeriesFilters specifies additional regular expressions to be applied on
	// the series names returned from the query. This is useful for constraints
	// that can't be represented in the SeriesQuery (e.g. series matching `container_.+`
	// not matching `container_.+_total`. A filter will be automatically appended to
	// match the form specified in Name.
	SeriesFilters []RegexFilter `json:"seriesFilters,omitempty"`
	// Resources specifies how associated Kubernetes resources should be discovered for
	// the given metrics.
	Resources ResourceMapping `json:"resources"`
	// Name specifies how the metric name should be transformed between custom metric
	// API resources, and Prometheus metric names.
	Name NameMapping `json:"name"`
	// MetricsQuery specifies modifications to the metrics query, such as converting
	// cumulative metrics to rate metrics. It is a template where `.LabelMatchers` is
	// a the comma-separated base label matchers and `.Series` is the series name, and
	// `.GroupBy` is the comma-separated expected group-by label names. The delimeters
	// are `<<` and `>>`.
	MetricsQuery string `json:"metricsQuery,omitempty"`
}

// RegexFilter is a filter that matches positively or negatively against a regex.
// Only one field may be set at a time.
type RegexFilter struct {
	Is    string `json:"is,omitempty"`
	IsNot string `json:"isNot,omitempty"`
}

// ResourceMapping specifies how to map Kubernetes resources to Prometheus labels
type ResourceMapping struct {
	// Template specifies a golang string template for converting a Kubernetes
	// group-resource to a Prometheus label.  The template object contains
	// the `.Group` and `.Resource` fields.  The `.Group` field will have
	// dots replaced with underscores, and the `.Resource` field will be
	// singularized.  The delimiters are `<<` and `>>`.
	Template string `json:"template,omitempty"`
	// Overrides specifies exceptions to the above template, mapping label names
	// to group-resources
	Overrides map[string]GroupResource `json:"overrides,omitempty"`
}

// GroupResource represents a Kubernetes group-resource.
type GroupResource struct {
	Group    string `json:"group,omitempty"`
	Resource string `json:"resource"`
}

// NameMapping specifies how to convert Prometheus metrics
// to/from custom metrics API resources.
type NameMapping struct {
	// Matches is a regular expression that is used to match
	// Prometheus series names.  It may be left blank, in which
	// case it is equivalent to `.*`.
	Matches string `json:"matches"`
	// As is the name used in the API.  Captures from Matches
	// are available for use here.  If not specified, it defaults
	// to $0 if no capture groups are present in Matches, or $1
	// if only one is present, and will error if multiple are.
	As string `json:"as,omitempty"`
}

// CustomScaleMetricRuleStatus defines the observed state of CustomScaleMetricRule
type CustomScaleMetricRuleStatus struct {
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// CustomScaleMetricRule is the Schema for the customscalemetricrules API
type CustomScaleMetricRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CustomScaleMetricRuleSpec   `json:"spec,omitempty"`
	Status CustomScaleMetricRuleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// CustomScaleMetricRuleList contains a list of CustomScaleMetricRule
type CustomScaleMetricRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CustomScaleMetricRule `json:"items"`
}

// ItemsToString returns items as <namespace>/<name> comma-separated string.
func (in *CustomScaleMetricRuleList) ItemsToString() string {
	if len(in.Items) == 0 {
		return ""
	}

	res := fmt.Sprintf("%s/%s", in.Items[0].Namespace, in.Items[0].Name)
	for _, itm := range in.Items[1:] {
		res += fmt.Sprintf(", %s/%s", itm.Namespace, itm.Name)
	}

	return res
}

func init() {
	SchemeBuilder.Register(&CustomScaleMetricRule{}, &CustomScaleMetricRuleList{})
}

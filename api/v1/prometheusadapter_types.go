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
	apiv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PrometheusAdapterSpec defines the desired state of PrometheusAdapter
type PrometheusAdapterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make generate" to regenerate code after modifying this file

	// Image to use for a `prometheus-adapter` deployment.
	Image string `json:"image"`

	// Replicas set the expected replicas of the prometheus-adapter. The controller will eventually make the size
	// of the running replicas equal to the expected size.
	Replicas *int32 `json:"replicas,omitempty"`

	// PrometheusURL used to connect to Prometheus. It will eventually contain query parameters
	// to configure the connection.
	PrometheusURL string `json:"prometheusUrl,omitempty"`

	// MetricsRelistInterval is the interval at which to update the cache of available metrics from Prometheus
	MetricsRelistInterval string `json:"metricsRelistInterval,omitempty"`

	// EnableResourceMetrics allows enabling/disabling adapter for `metrics.k8s.io`
	EnableResourceMetrics bool `json:"enableResourceMetrics,omitempty"`

	// EnableCustomMetrics allows enabling/disabling adapter for `custom.metrics.k8s.io`
	EnableCustomMetrics bool `json:"enableCustomMetrics,omitempty"`

	// CustomScaleMetricRulesSelector defines label selectors to select
	// CustomScaleMetricRule resources across the cluster.
	CustomScaleMetricRulesSelector []*metav1.LabelSelector `json:"customScaleMetricRulesSelector,omitempty"`

	// Resources defines resources requests and limits for single Pods.
	Resources v1.ResourceRequirements `json:"resources,omitempty"`

	// SecurityContext holds pod-level security attributes.
	SecurityContext *SecurityContext `json:"securityContext,omitempty"`

	TLSEnabled bool `json:"tlsEnabled,omitempty"`

	TLSConfig *TLSConfig `json:"tlsConfig,omitempty"`

	Auth *Auth `json:"auth,omitempty"`

	// Define which Nodes the Pods are scheduled on.
	// Specified just as map[string]string. For example: "type: compute"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) objects. May match selectors of replication controllers
	// and services.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations is an unstructured key value map stored with a resource that may be
	// set by external tools to store and retrieve arbitrary metadata. They are not
	// queryable and should be preserved when modifying objects.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Affinity is a group of affinity scheduling rules.
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node
	// +optional
	Affinity *v1.Affinity `json:"affinity,omitempty"`

	// Tolerations allow the pods to schedule onto nodes with matching taints.
	// More info: https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration
	// +optional
	Tolerations []v1.Toleration `json:"tolerations,omitempty"`

	// PriorityClassName assigned to the Pods
	// +optional
	PriorityClassName string `json:"priorityClassName,omitempty"`
}

// PrometheusAdapterStatus defines the observed state of PrometheusAdapter
type PrometheusAdapterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// PrometheusAdapter is the Schema for the prometheusadapters API
type PrometheusAdapter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrometheusAdapterSpec   `json:"spec,omitempty"`
	Status PrometheusAdapterStatus `json:"status,omitempty"`
}

// SecurityContext holds pod-level security attributes.
// The parameters are required if a Pod Security Policy is enabled
// for Kubernetes cluster and required if a Security Context Constraints is enabled
// for Openshift cluster.
type SecurityContext struct {
	// The UID to run the entrypoint of the container process.
	// Defaults to user specified in image metadata if unspecified.
	RunAsUser *int64 `json:"runAsUser,omitempty"`
	// A special supplemental group that applies to all containers in a pod.
	// Some volume types allow the Kubelet to change the ownership of that volume
	// to be owned by the pod:
	//
	// 1. The owning GID will be the FSGroup
	// 2. The setgid bit is set (new files created in the volume will be owned by FSGroup)
	// 3. The permission bits are OR'd with rw-rw----
	//
	// If unset, the Kubelet will not modify the ownership and permissions of any volume.
	FSGroup *int64 `json:"fsGroup,omitempty"`
}

// +kubebuilder:object:root=true

// PrometheusAdapterList contains a list of PrometheusAdapter
type PrometheusAdapterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrometheusAdapter `json:"items"`
}

type TLSConfig struct {
	// Certificate authority used when verifying server certificates.
	CA *v1.SecretKeySelector `json:"caSecret,omitempty"`
	// Client certificate to present when doing client-authentication.
	Cert *v1.SecretKeySelector `json:"certSecret,omitempty"`
	// Secret containing the client key file for the target.
	KeySecret *v1.SecretKeySelector `json:"keySecret,omitempty"`
}

type Auth struct {
	BasicAuth *apiv1.BasicAuth `json:"basicAuth,omitempty"`
}

func init() {
	SchemeBuilder.Register(&PrometheusAdapter{}, &PrometheusAdapterList{})
}

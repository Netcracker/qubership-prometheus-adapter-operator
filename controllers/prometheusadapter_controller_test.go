// Copyright NetCracker Technology Corporation
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

package controllers

import (
	"testing"
)

func TestGetTagFromImage(t *testing.T) {
	tests := []struct {
		name     string
		image    string
		expected string
	}{
		{
			name:     "image with tag",
			image:    "prometheus-adapter:v0.10.0",
			expected: "v0.10.0",
		},
		{
			name:     "image with registry and tag",
			image:    "quay.io/prometheus-adapter/prometheus-adapter:v0.10.0",
			expected: "v0.10.0",
		},
		{
			name:     "image without tag",
			image:    "prometheus-adapter",
			expected: "prometheus-adapter",
		},
		{
			name:     "image with multiple colons",
			image:    "registry:5000/prometheus-adapter:v0.10.0",
			expected: "v0.10.0",
		},
		{
			name:     "empty string",
			image:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTagFromImage(tt.image)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetInstanceLabel(t *testing.T) {
	tests := []struct {
		name      string
		labelName string
		namespace string
		expected  string
	}{
		{
			name:      "normal case",
			labelName: "prometheus-adapter",
			namespace: "monitoring",
			expected:  "prometheus-adapter-monitoring",
		},
		{
			name:      "long name that exceeds 63 chars",
			labelName: "very-long-prometheus-adapter-name-that-exceeds-sixty-three-characters",
			namespace: "very-long-namespace-name-that-also-exceeds-limits",
			expected:  "very-long-prometheus-adapter-name-that-exceeds-sixty-three-char",
		},
		{
			name:      "name with trailing dash",
			labelName: "test-",
			namespace: "ns",
			expected:  "test--ns",
		},
		{
			name:      "empty name",
			labelName: "",
			namespace: "ns",
			expected:  "ns",
		},
		{
			name:      "empty namespace",
			labelName: "test",
			namespace: "",
			expected:  "test",
		},
		{
			name:      "both empty",
			labelName: "",
			namespace: "",
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getInstanceLabel(tt.labelName, tt.namespace)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
			if len(result) > 63 {
				t.Errorf("result length %d exceeds 63 characters", len(result))
			}
		})
	}
}

func TestCheckDeploymentStatus(t *testing.T) {
	// Note: This test would require setting up a fake client and deployment
	// For now, we'll test the logic that can be tested without external dependencies
	// The checkDeploymentStatus method is complex and depends on k8s client calls,
	// so it's better suited for integration tests
	t.Skip("checkDeploymentStatus requires k8s client setup, better for integration tests")
}

// Test that the Factory methods create proper manifests
// These would be integration tests since they test the manifest creation
func TestFactoryMethods(t *testing.T) {
	t.Skip("Factory methods create k8s manifests, better tested as integration tests")
}

func TestPrometheusAdapterReconciler_SetupWithManager(t *testing.T) {
	t.Skip("SetupWithManager requires a real manager, better for integration tests")
}

// Test the Reconcile method logic that can be tested without k8s dependencies
func TestReconcileLogic(t *testing.T) {
	t.Skip("Reconcile logic is complex and requires mocking, better for integration tests")
}

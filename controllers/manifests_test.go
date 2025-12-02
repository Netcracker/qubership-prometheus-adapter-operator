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

package controllers_test

import (
	"bufio"
	"os"
	"testing"

	api "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"

	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var (
	cr              *api.PrometheusAdapter
	factory         *controllers.Factory
	labelKey        = "label.key"
	labelValue      = "label-value"
	annotationKey   = "annotation.key"
	annotationValue = "annotation-value"
)

func TestPrometheusAdapterDeploymentManifest(t *testing.T) {
	d := appsv1.Deployment{}
	err := yaml.NewYAMLOrJSONDecoder(controllers.MustAssetReader(controllers.PrometheusAdapterDeployment), 100).Decode(&d)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, d, "Deployment manifest should not be empty")

	dep := appsv1.Deployment{}
	f, err := os.Open(controllers.PrometheusAdapterDeployment)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&dep)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, dep, "Deployment manifest should not be empty")

	assert.Equal(t, dep, d)
}

func TestPrometheusAdapterServiceManifest(t *testing.T) {
	s := corev1.Service{}
	err := yaml.NewYAMLOrJSONDecoder(controllers.MustAssetReader(controllers.PrometheusAdapterService), 100).Decode(&s)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s, "Service manifest should not be empty")

	srv := corev1.Service{}
	f, err := os.Open(controllers.PrometheusAdapterService)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&srv)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, srv, "Service manifest should not be empty")

	assert.Equal(t, srv, s)
}

func TestPrometheusAdapterServiceAccountManifest(t *testing.T) {
	s := corev1.ServiceAccount{}
	err := yaml.NewYAMLOrJSONDecoder(controllers.MustAssetReader(controllers.PrometheusAdapterServiceAccount), 100).Decode(&s)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, s, "Service account manifest should not be empty")

	sa := corev1.ServiceAccount{}
	f, err := os.Open(controllers.PrometheusAdapterServiceAccount)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&sa)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, sa, "Service account manifest should not be empty")

	assert.Equal(t, sa, s)
}

func TestCustomMetricsClusterRoleManifest(t *testing.T) {
	clusterRole := rbacv1.ClusterRole{}
	err := yaml.NewYAMLOrJSONDecoder(controllers.MustAssetReader(controllers.CustomMetricsClusterRole), 100).Decode(&clusterRole)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, clusterRole, "Cluster Role manifest should not be empty")

	cr := rbacv1.ClusterRole{}
	f, err := os.Open(controllers.CustomMetricsClusterRole)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&cr)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, cr, "Cluster Role manifest should not be empty")

	assert.Equal(t, cr, clusterRole)
}

func TestCustomMetricsClusterRoleBindingManifest(t *testing.T) {
	clusterRoleBinding := rbacv1.ClusterRoleBinding{}
	err := yaml.NewYAMLOrJSONDecoder(controllers.MustAssetReader(controllers.CustomMetricsClusterRoleBinding), 100).Decode(&clusterRoleBinding)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, clusterRoleBinding, "Cluster Role Binding manifest should not be empty")

	crb := rbacv1.ClusterRoleBinding{}
	f, err := os.Open(controllers.CustomMetricsClusterRoleBinding)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&crb)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, crb, "Cluster Role Binding manifest should not be empty")

	assert.Equal(t, crb, clusterRoleBinding)
}

func TestCustomMetricsConfigMapManifest(t *testing.T) {
	configMap := corev1.ConfigMap{}
	err := yaml.NewYAMLOrJSONDecoder(controllers.MustAssetReader(controllers.PrometheusAdapterConfigMap), 100).Decode(&configMap)

	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, configMap, "Config Map manifest should not be empty")

	cm := corev1.ConfigMap{}
	f, err := os.Open(controllers.PrometheusAdapterConfigMap)
	if err != nil {
		t.Fatal(err)
	}
	err = yaml.NewYAMLOrJSONDecoder(bufio.NewReader(f), 100).Decode(&cm)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotNil(t, cm, "Config Map manifest should not be empty")

	assert.Equal(t, cm, configMap)
}

func TestPrometheusAdapterOperatorManifests(t *testing.T) {
	cr = &api.PrometheusAdapter{
		Spec: api.PrometheusAdapterSpec{
			Annotations: map[string]string{annotationKey: annotationValue},
			Labels:      map[string]string{labelKey: labelValue},
		},
	}
	factory = controllers.NewFactory(cr)

	t.Run("Test Deployment manifest", func(t *testing.T) {
		m, err := factory.PrometheusAdapterDeployment()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Equal(t, labelValue, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Equal(t, labelValue, m.Spec.Template.Labels[labelKey])
		assert.NotNil(t, m.GetAnnotations())
		assert.Equal(t, annotationValue, m.GetAnnotations()[annotationKey])
		assert.Equal(t, annotationValue, m.Spec.Template.Annotations[annotationKey])
	})
	cr = &api.PrometheusAdapter{
		Spec: api.PrometheusAdapterSpec{},
	}
	factory = controllers.NewFactory(cr)
	t.Run("Test Deployment manifest with nil labels and annotation", func(t *testing.T) {
		m, err := factory.PrometheusAdapterDeployment()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Deployment manifest should not be empty")
		assert.NotNil(t, m.GetLabels())
		assert.Empty(t, m.GetLabels()[labelKey])
		assert.NotNil(t, m.Spec.Template.Labels)
		assert.Empty(t, m.Spec.Template.Labels[labelKey])
		assert.Nil(t, m.GetAnnotations())
		assert.Empty(t, m.Spec.Template.Annotations[annotationKey])
	})
	t.Run("Test Service manifest", func(t *testing.T) {
		m, err := factory.PrometheusAdapterService()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service manifest should not be empty")
	})
	t.Run("Test Service Account manifest", func(t *testing.T) {
		m, err := factory.PrometheusAdapterServiceAccount()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Service Account manifest should not be empty")
	})
	t.Run("Test Custom Metric Cluster Role manifest", func(t *testing.T) {
		m, err := factory.CustomMetricsClusterRole()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Custom Metric Cluster Role manifest should not be empty")
	})
	t.Run("Test Custom Metric Cluster Role Binding manifest", func(t *testing.T) {
		m, err := factory.CustomMetricsClusterRoleBinding()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Custom Metric Cluster Role Binding manifest should not be empty")
	})
	t.Run("Test Custom Metric Config Map manifest", func(t *testing.T) {
		m, err := factory.CustomMetricsConfigMap()
		if err != nil {
			t.Fatal(err)
		}
		assert.NotNil(t, m, "Custom Metric Config Map manifest should not be empty")
	})
}

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
	"context"
	"fmt"
	"time"

	promv1 "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers/config"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/yaml"
)

// RebuildPrometheusAdapterConfig rebuilds and updates prometheus-adapter configmap.
// Function flow:
//
//	# Step 1: Receive all CustomScaleMetricRule-s according to label selectors
//	# Step 2: Exclude (if necessary) deleted rules
//	# Step 3: Convert actual lists of rules to JSON
//	# Step 4: Update prometheus-adapter ConfigMap
func (m *PrometheusAdapterManager) RebuildPrometheusAdapterConfig() error {
	configMapModel := &ConfigMap{
		CustomMetricRules: []promv1.CustomMetricRuleConfig{},
		ResourceRules:     ResourceRules{},
	}

	m.log.Info("updating rules: start")
	cfg := config.GetControllerConfig()
	if !cfg.GetEnableResourceMetrics() && !cfg.GetEnableCustomMetrics() {
		m.log.Info("resourceMetrics and customMetrics disabled for prometheus-adapter")
		return nil
	}
	if cfg.GetEnableResourceMetrics() {
		configMap := corev1.ConfigMap{}
		err := m.client.Get(context.TODO(), apimachinerytypes.NamespacedName{Name: "prometheus-adapter-resource-rules", Namespace: cfg.GetActivatedBy().Namespace}, &configMap)
		if errors.IsNotFound(err) {
			m.log.Info("failed to get prometheus-adapter-resource-rules configmap. use default configmap")
			configMapModel.ResourceRules = defaultResourceRules
		} else {
			err = yaml.Unmarshal([]byte(configMap.Data["config.yaml"]), &configMapModel)
			if err != nil {
				m.log.Error(err, "Error unmarshal configMap content")
			}
		}
		configMapModel.CustomMetricRules = []promv1.CustomMetricRuleConfig{}
	}
	if cfg.GetEnableCustomMetrics() {
		allCustomMetricRules := &promv1.CustomScaleMetricRuleList{}
		// Build label selectors as ListOptions
		var options []client.ListOption
		for _, ls := range cfg.GetCustomMetricRulesSelectors() {
			s, _ := metav1.LabelSelectorAsSelector(ls)
			options = append(options, client.MatchingLabelsSelector{Selector: s})
		}
		if err := m.client.List(context.TODO(), allCustomMetricRules, options...); err != nil {
			m.log.Error(err, "updating custom rules: failed get list of resources")
			return err
		}

		m.log.Info(fmt.Sprintf("updating custom rules: actual custom metric rules: %v", allCustomMetricRules.ItemsToString()))
		// Marshall all rules to JSON
		for _, rule := range allCustomMetricRules.Items {
			configMapModel.CustomMetricRules = append(configMapModel.CustomMetricRules, rule.Spec.Rules...)
		}
	}
	configMap := corev1.ConfigMap{}
	err := m.client.Get(context.TODO(), apimachinerytypes.NamespacedName{Name: "prometheus-adapter-config", Namespace: cfg.GetActivatedBy().Namespace}, &configMap)
	if err != nil {
		m.log.Error(err, "updating custom rules: failed get prometheus-adapter configmap")
		return err
	}

	rawConfigMap, err := yaml.Marshal(configMapModel)
	if err != nil {
		m.log.Error(err, "updating custom rules: failed marshall prometheus-adapter ConfigMap with metrics rules")
		return err
	}

	// Update config map
	lockTmt := time.Second
	err = cfg.LockConfigMap(&lockTmt)
	if err != nil {
		return err
	}
	defer cfg.UnlockConfigMap()

	configMap.Data["config.yaml"] = string(rawConfigMap)
	m.log.Info("update configmap")
	err = m.client.Update(context.TODO(), &configMap)
	if err != nil {
		m.log.Error(err, "updating rules: failed update prometheus-adapter configmap")
		return err
	}

	return nil
}

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

package config

import (
	"testing"
	"time"

	apimachinerymetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
)

func TestGetControllerConfig(t *testing.T) {
	// Test singleton pattern
	cfg1 := GetControllerConfig()
	cfg2 := GetControllerConfig()

	if cfg1 != cfg2 {
		t.Error("GetControllerConfig should return the same instance")
	}
}

func TestControllerConfig_IsActivated(t *testing.T) {
	cfg := GetControllerConfig()
	cfg.Deactivate() // Reset to known state

	if cfg.IsActivated() {
		t.Error("expected config to be deactivated initially")
	}

	cfg.Activate()
	if !cfg.IsActivated() {
		t.Error("expected config to be activated after calling Activate")
	}

	cfg.Deactivate()
	if cfg.IsActivated() {
		t.Error("expected config to be deactivated after calling Deactivate")
	}
}

func TestControllerConfig_GetActivatedBy(t *testing.T) {
	cfg := GetControllerConfig()
	cfg.Deactivate() // Reset to known state

	// Initially should be nil
	if cfg.GetActivatedBy() != nil {
		t.Error("expected activatedBy to be nil initially")
	}

	// Set activated by
	namespacedName := &apimachinerytypes.NamespacedName{
		Name:      "test-name",
		Namespace: "test-namespace",
	}
	cfg.SetActivatedBy(namespacedName)

	if cfg.GetActivatedBy() == nil {
		t.Error("expected activatedBy to be set")
	}

	if cfg.GetActivatedBy().Name != "test-name" || cfg.GetActivatedBy().Namespace != "test-namespace" {
		t.Error("expected activatedBy to match set value")
	}
}

func TestControllerConfig_CustomMetricRulesSelectors(t *testing.T) {
	cfg := GetControllerConfig()

	// Initially should be empty label selector
	selectors := cfg.GetCustomMetricRulesSelectors()
	if len(selectors) != 1 || len(selectors[0].MatchLabels) != 0 {
		t.Error("expected initial selectors to be empty label selector")
	}

	// Set custom selectors
	customSelectors := []*apimachinerymetav1.LabelSelector{
		{MatchLabels: map[string]string{"app": "test"}},
	}
	cfg.SetCustomMetricRulesSelectors(customSelectors)

	result := cfg.GetCustomMetricRulesSelectors()
	if len(result) != 1 || result[0].MatchLabels["app"] != "test" {
		t.Error("expected custom selectors to be set correctly")
	}
}

func TestControllerConfig_EnableAdapters(t *testing.T) {
	cfg := GetControllerConfig()

	// Set enabled adapters
	cfg.SetEnabledAdapters(true, false)

	if !cfg.GetEnableResourceMetrics() {
		t.Error("expected resource metrics to be enabled")
	}

	if cfg.GetEnableCustomMetrics() {
		t.Error("expected custom metrics to be disabled")
	}

	// Test other combination
	cfg.SetEnabledAdapters(false, true)

	if cfg.GetEnableResourceMetrics() {
		t.Error("expected resource metrics to be disabled")
	}

	if !cfg.GetEnableCustomMetrics() {
		t.Error("expected custom metrics to be enabled")
	}
}

func TestControllerConfig_LockConfigMap(t *testing.T) {
	cfg := GetControllerConfig()

	// Initially should not be locked
	if cfg.IfConfigMapLocked() {
		t.Error("expected config map to not be locked initially")
	}

	// Lock without timeout
	err := cfg.LockConfigMap(nil)
	if err != nil {
		t.Errorf("expected no error when locking initially, got: %v", err)
	}

	if !cfg.IfConfigMapLocked() {
		t.Error("expected config map to be locked")
	}

	// Try to lock again without timeout (should fail)
	err = cfg.LockConfigMap(nil)
	if err == nil {
		t.Error("expected error when locking already locked config map")
	}

	// Unlock
	cfg.UnlockConfigMap()
	if cfg.IfConfigMapLocked() {
		t.Error("expected config map to be unlocked")
	}

	// Test with timeout
	timeout := 1 * time.Second
	err = cfg.LockConfigMap(&timeout)
	if err != nil {
		t.Errorf("expected no error when locking with timeout, got: %v", err)
	}

	if !cfg.IfConfigMapLocked() {
		t.Error("expected config map to be locked")
	}

	cfg.UnlockConfigMap()
}

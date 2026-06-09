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

package config

import (
	"errors"
	"fmt"
	"sync"
	"time"

	apimachinerymetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"

	klog "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	// EmptyLabelSelector is a label selector which matches to any object.
	EmptyLabelSelector = []*apimachinerymetav1.LabelSelector{{}}
)

// ControllerConfig defines thread safe settings for qubership-prometheus-adapter-operator controller.
// This struct works is a singleton to share values between reconciliation cycles for
// different controllers.
type ControllerConfig struct {
	*sync.Mutex

	// activated marks has we already reconciled custom resource PrometheusAdapter.
	// It is simple restriction for having only one prometheus-adapter reconciled by
	// an instance of qubership-prometheus-adapter-operator. Several instances of
	// prometheus-adapter which handle several API can not work fine in one K8s cluster.
	activated bool
	// activatedBy stores namespace and names of the first reconciled PrometheusAdapter custom resource.
	activatedBy *apimachinerytypes.NamespacedName

	// customMetricRulesSelectors stores a set of label selectors to watch CustomScaleMetricRule custom resources.
	customMetricRulesSelectors []*apimachinerymetav1.LabelSelector

	// configMapLocked indicates that prometheus-adapter configmap is updating now
	configMapLocked bool

	// enableResourceMetrics indicates that prometheus-adapter can work with `metrics.k8s.io`
	enableResourceMetrics bool

	// enableCustomMetrics indicates that prometheus-adapter can work with `custom.metrics.k8s.io`
	enableCustomMetrics bool
}

var instance *ControllerConfig
var once sync.Once

var log = klog.Log.WithName("controller-config")

// GetControllerConfig creates an instance of ControllerConfig or returns existed one.
func GetControllerConfig() *ControllerConfig {
	once.Do(func() {
		instance = &ControllerConfig{
			Mutex:                      &sync.Mutex{},
			activated:                  false,
			activatedBy:                nil,
			customMetricRulesSelectors: EmptyLabelSelector,
		}
	})

	return instance
}

// IsActivated returns value of activated property
func (c *ControllerConfig) IsActivated() bool {
	return c.activated
}

// Activate sets value of activated property as true
func (c *ControllerConfig) Activate() {
	c.Lock()
	defer c.Unlock()
	c.activated = true
}

// Deactivate sets value of activated property as false
func (c *ControllerConfig) Deactivate() {
	c.Lock()
	defer c.Unlock()
	c.activated = false
}

// GetActivatedBy returns value of activatedBy property
func (c *ControllerConfig) GetActivatedBy() *apimachinerytypes.NamespacedName {
	return c.activatedBy
}

// SetActivatedBy sets activatedBy property
func (c *ControllerConfig) SetActivatedBy(v *apimachinerytypes.NamespacedName) {
	c.Lock()
	defer c.Unlock()
	c.activatedBy = v
}

// GetcustomMetricRulesSelectors returns value of customMetricRulesSelectors property
func (c *ControllerConfig) GetCustomMetricRulesSelectors() []*apimachinerymetav1.LabelSelector {
	c.Lock()
	defer c.Unlock()
	return c.customMetricRulesSelectors
}

// SetcustomMetricRulesSelectors sets customMetricRulesSelectors property
func (c *ControllerConfig) SetCustomMetricRulesSelectors(v []*apimachinerymetav1.LabelSelector) {
	c.Lock()
	defer c.Unlock()
	c.customMetricRulesSelectors = v
}

// SetEnabledAdapters set enableResourceMetrics and enableCustomMetrics properties
func (c *ControllerConfig) SetEnabledAdapters(resourceMetrics, customMetrics bool) {
	c.Lock()
	defer c.Unlock()
	c.enableResourceMetrics = resourceMetrics
	c.enableCustomMetrics = customMetrics
}

// GetEnableResourceMetrics returns value of enableResourceMetrics property
func (c *ControllerConfig) GetEnableResourceMetrics() bool {
	c.Lock()
	defer c.Unlock()
	return c.enableResourceMetrics
}

// GetEnableCustomMetrics returns value of enableCustomMetrics property
func (c *ControllerConfig) GetEnableCustomMetrics() bool {
	c.Lock()
	defer c.Unlock()
	return c.enableCustomMetrics
}

// IfConfigMapLocked returns value of configMapLocked property.
func (c *ControllerConfig) IfConfigMapLocked() bool {
	return c.configMapLocked
}

func (c *ControllerConfig) UnlockConfigMap() {
	c.Lock()
	defer c.Unlock()
	c.configMapLocked = false
	log.Info("unlock ConfigMap")
}

// LockConfigMap sets configMapLocked as true with timeout.
func (c *ControllerConfig) LockConfigMap(timeout *time.Duration) error {
	if timeout == nil {
		c.Lock()
		defer c.Unlock()
		if !c.configMapLocked {
			c.configMapLocked = true
			log.Info("lock ConfigMap")
			return nil
		}

		return errors.New("ConfigMap is already locked")
	}

	success := false
	delay := 1 * time.Second
	attempts := 1
	if *timeout > delay {
		attempts = int(*timeout / delay)
	}
	for i := 0; i < attempts; i++ {
		if !c.configMapLocked {
			c.Lock()
			c.configMapLocked = true
			c.Unlock()
			success = true
			log.Info("lock ConfigMap")
			break
		} else if i+1 < attempts {
			log.Info("waiting ConfigMap to lock")
			time.Sleep(delay)
		}
	}

	if success {
		return nil
	} else {
		return fmt.Errorf("ConfigMap is already locked, can not allocate it within %s", timeout.String())

	}
}

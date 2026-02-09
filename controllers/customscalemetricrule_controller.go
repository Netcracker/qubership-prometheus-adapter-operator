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

package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"

	promv1 "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers/common"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers/config"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers/prometheusadapter"
)

// CustomScaleMetricRuleReconciler reconciles a CustomScaleMetricRule object
type CustomScaleMetricRuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconciles a CustomScaleMetricRule object.
func (r *CustomScaleMetricRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("customscalemetricrule", req.NamespacedName)
	cfg := config.GetControllerConfig()

	// Skip reconciliation if prometheus-adapter is not reconciled yet and requeue
	if !cfg.IsActivated() || cfg.GetActivatedBy() == nil {
		log.Info(fmt.Sprintf("skip reconciliation: there is no reconciled prometheus-adapter (retry after %s)", OneMinute.String()))
		return ctrl.Result{RequeueAfter: OneMinute}, nil
	}

	log.Info("start reconcile")

	prometheusAdapterManager := prometheusadapter.NewPrometheusAdapterManager(r.Client, r.Log)

	// Get CR instance
	customScaleMetricRule := &promv1.CustomScaleMetricRule{}
	err := r.Get(context.TODO(), req.NamespacedName, customScaleMetricRule)
	if err != nil {
		if errors.IsNotFound(err) {
			if err := prometheusAdapterManager.RebuildPrometheusAdapterConfig(); err != nil {
				return ctrl.Result{}, err
			}

			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizes.
			// Return and don't requeue
			return ctrl.Result{Requeue: false}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Check label selectors
	match, err := common.MatchAll(customScaleMetricRule.GetObjectMeta(), cfg.GetCustomMetricRulesSelectors())
	if err != nil {
		log.Error(err, "error matching selectors")
		return ctrl.Result{RequeueAfter: 60 * time.Second}, nil
	} else if !match {
		return ctrl.Result{Requeue: false}, nil
	}

	if err := prometheusAdapterManager.RebuildPrometheusAdapterConfig(); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager specifies how the controller is built to watch a CR and other resources that are owned and managed by that controller.
func (r *CustomScaleMetricRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&promv1.CustomScaleMetricRule{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 1,
		}).
		Complete(r)
}

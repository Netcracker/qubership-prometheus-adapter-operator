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
	"reflect"
	"strings"
	"time"

	promv1 "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers/config"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers/prometheusadapter"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PrometheusAdapterReconciler reconciles a PrometheusAdapter object
type PrometheusAdapterReconciler struct {
	Client client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// Reconcile reconciles a PrometheusAdapter object
func (r *PrometheusAdapterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	cfg := config.GetControllerConfig()
	log := r.Log.WithValues("prometheusadapter", req.NamespacedName)
	log.Info("start reconcile")

	customResourceInstance := &promv1.PrometheusAdapter{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, customResourceInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("deactivate prometheus-adapter")
			cfg.Deactivate()
			cfg.SetActivatedBy(nil)
			cfg.SetCustomMetricRulesSelectors(config.EmptyLabelSelector)
			cfg.SetEnabledAdapters(customResourceInstance.Spec.EnableResourceMetrics, customResourceInstance.Spec.EnableCustomMetrics)
			return ctrl.Result{Requeue: false}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	cfg.SetEnabledAdapters(customResourceInstance.Spec.EnableResourceMetrics, customResourceInstance.Spec.EnableCustomMetrics)

	// Check if we can reconcile this resource
	if cfg.IsActivated() && cfg.GetActivatedBy().Name != req.Name && cfg.GetActivatedBy().Namespace != req.Namespace {
		log.Info(fmt.Sprintf("skip reconciliation: now reconcile %s (retry after %s)", cfg.GetActivatedBy().String(), ThreeMinutes.String()))
		return ctrl.Result{RequeueAfter: ThreeMinutes}, nil
	}

	if !cfg.IsActivated() {
		log.Info("activate prometheus-adapter")
		cfg.Activate()
		cfg.SetActivatedBy(&req.NamespacedName)
	}

	// Manage prometheus-adapter ConfigMap
	if err := r.handleConfigMap(customResourceInstance); err != nil {
		return ctrl.Result{}, err
	}

	// Manage prometheus-adapter ServiceAccount
	if err := r.handleServiceAccount(customResourceInstance); err != nil {
		return ctrl.Result{}, err
	}

	// Manage prometheus-adapter ClusterRole
	if err := r.handleClusterRole(customResourceInstance); err != nil {
		return ctrl.Result{}, err
	}

	// Manage prometheus-adapter ClusterRoleBinding
	if err := r.handleClusterRoleBinding(customResourceInstance); err != nil {
		return ctrl.Result{}, err
	}

	// Manage prometheus-adapter Service
	if err := r.handleService(customResourceInstance); err != nil {
		return ctrl.Result{}, err
	}

	// Manage prometheus-adapter Deployment
	if err := r.handleDeployment(customResourceInstance); err != nil {
		return ctrl.Result{}, err
	}

	isDeploymentAvailable := r.checkDeploymentStatus(customResourceInstance)
	if !isDeploymentAvailable {
		log.Info("Deployment is unavailable. Requeueing reconciliation")
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// Fill controllers config and update configmap
	oldSelectors := cfg.GetCustomMetricRulesSelectors()

	if len(customResourceInstance.Spec.CustomScaleMetricRulesSelector) > 0 {
		cfg.SetCustomMetricRulesSelectors(customResourceInstance.Spec.CustomScaleMetricRulesSelector)
	} else {
		cfg.SetCustomMetricRulesSelectors(config.EmptyLabelSelector)
	}

	if !reflect.DeepEqual(oldSelectors, cfg.GetCustomMetricRulesSelectors()) {
		prometheusAdapterManager := prometheusadapter.NewPrometheusAdapterManager(r.Client, r.Log)
		if err := prometheusAdapterManager.RebuildPrometheusAdapterConfig(); err != nil {
			return ctrl.Result{}, err
		}
	}

	log.Info("reconciliation finished")
	return ctrl.Result{}, nil
}

// SetupWithManager creates a contreller for PrometheusAdapter
func (r *PrometheusAdapterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&promv1.PrometheusAdapter{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *PrometheusAdapterReconciler) createResource(cr *promv1.PrometheusAdapter, o K8sResource) error {
	oMeta := GetKindedNamespacedName(o)
	log := r.Log.WithValues(
		"prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()),
		"resourceKind", oMeta.Kind,
		"resource", oMeta.NamespacedName(),
	)

	if err := controllerutil.SetControllerReference(cr, o, r.Scheme); err != nil {
		return err
	}

	if err := r.Client.Create(context.TODO(), o); err != nil {
		if errors.IsAlreadyExists(err) {
			return err
		}

		log.Error(err, "resource creating: resource create failed")
		return err
	}

	log.Info("resource creating: resource created")
	return nil
}

func (r *PrometheusAdapterReconciler) updateResource(cr *promv1.PrometheusAdapter, o K8sResource) error {
	oMeta := GetKindedNamespacedName(o)
	log := r.Log.WithValues(
		"prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()),
		"resourceKind", oMeta.Kind,
		"resource", oMeta.NamespacedName(),
	)

	// Update object
	if err := r.Client.Update(context.TODO(), o); err != nil {
		log.Error(err, "resource updating: resource update failed")
		return err
	}

	log.Info("resource updating: resource updated")
	return nil
}

// getResource tries to get resource inside namespace or on cluster level.
func (r *PrometheusAdapterReconciler) getResource(o K8sResource) error {
	objectKey := client.ObjectKeyFromObject(o)
	if err := r.Client.Get(context.TODO(), objectKey, o); err != nil {
		if errors.IsNotFound(err) {
			objectKey.Namespace = ""
			if err := r.Client.Get(context.TODO(), objectKey, o); err == nil {
				return nil
			}
		}
		return err
	}
	return nil
}

// checkDeploymentStatus check if all deployment pods are available
func (r *PrometheusAdapterReconciler) checkDeploymentStatus(cr *promv1.PrometheusAdapter) bool {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	d := &appsv1.Deployment{}
	deploymentName := "prometheus-adapter"

	err := r.Client.Get(
		context.TODO(),
		types.NamespacedName{
			Name:      deploymentName,
			Namespace: cr.GetNamespace(),
		},
		d,
	)

	if err != nil {
		log.Error(err, "failed Get deployment %s", deploymentName)
		return false
	}

	// TODO: What's the right way to check availability?
	return d.Status.UnavailableReplicas == 0 && d.Status.AvailableReplicas > 0
}

func (r *PrometheusAdapterReconciler) handleServiceAccount(cr *promv1.PrometheusAdapter) error {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	f := NewFactory(cr)
	m, err := f.PrometheusAdapterServiceAccount()
	if err != nil {
		log.Error(err, "failed creating ServiceAccount manifest")
		return err
	}

	// Set labels
	m.Labels["app.kubernetes.io/instance"] = getInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = getTagFromImage(cr.Spec.Image)

	if err := r.createResource(cr, m); err != nil {
		if errors.IsAlreadyExists(err) {
			e := &corev1.ServiceAccount{ObjectMeta: m.ObjectMeta}
			if err := r.getResource(e); err != nil {
				return err
			}

			e.SetLabels(m.GetLabels())

			if err := r.updateResource(cr, e); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PrometheusAdapterReconciler) handleClusterRole(cr *promv1.PrometheusAdapter) error {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	f := NewFactory(cr)
	m, err := f.CustomMetricsClusterRole()
	if err != nil {
		log.Error(err, "failed creating ClusterRole manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = m.GetName()
	m.Labels["app.kubernetes.io/name"] = m.GetName()
	m.Labels["app.kubernetes.io/instance"] = getInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = getTagFromImage(cr.Spec.Image)

	if err := r.createResource(cr, m); err != nil {
		if errors.IsAlreadyExists(err) {
			e := &rbacv1.ClusterRole{ObjectMeta: m.ObjectMeta}
			if err := r.getResource(e); err != nil {
				return err
			}

			e.SetLabels(m.GetLabels())
			e.Rules = m.Rules

			if err := r.updateResource(cr, e); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PrometheusAdapterReconciler) handleClusterRoleBinding(cr *promv1.PrometheusAdapter) error {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	f := NewFactory(cr)
	m, err := f.CustomMetricsClusterRoleBinding()
	if err != nil {
		log.Error(err, "failed creating ClusterRoleBinding manifest")
		return err
	}

	// Set labels
	m.Labels["name"] = m.GetName()
	m.Labels["app.kubernetes.io/name"] = m.GetName()
	m.Labels["app.kubernetes.io/instance"] = getInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = getTagFromImage(cr.Spec.Image)

	if err := r.createResource(cr, m); err != nil {
		if errors.IsAlreadyExists(err) {
			e := &rbacv1.ClusterRoleBinding{ObjectMeta: m.ObjectMeta}
			if err := r.getResource(e); err != nil {
				return err
			}

			e.SetLabels(m.GetLabels())
			e.Subjects = m.Subjects
			e.RoleRef = m.RoleRef

			if err := r.updateResource(cr, e); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PrometheusAdapterReconciler) handleConfigMap(cr *promv1.PrometheusAdapter) error {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	f := NewFactory(cr)
	m, err := f.CustomMetricsConfigMap()
	if err != nil {
		log.Error(err, "failed creating ConfigMap manifest")
		return err
	}

	// Set labels
	m.Labels["app.kubernetes.io/instance"] = getInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = getTagFromImage(cr.Spec.Image)

	if err := r.createResource(cr, m); err != nil {
		if errors.IsAlreadyExists(err) {
			e := &corev1.ConfigMap{ObjectMeta: m.ObjectMeta}
			if err := r.getResource(e); err != nil {
				return err
			}

			e.SetLabels(m.GetLabels())

			if err := r.updateResource(cr, e); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PrometheusAdapterReconciler) handleService(cr *promv1.PrometheusAdapter) error {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	f := NewFactory(cr)
	m, err := f.PrometheusAdapterService()
	if err != nil {
		log.Error(err, "failed creating Service manifest")
		return err
	}

	// Set labels
	m.Labels["app.kubernetes.io/instance"] = getInstanceLabel(m.GetName(), m.GetNamespace())
	m.Labels["app.kubernetes.io/version"] = getTagFromImage(cr.Spec.Image)

	if err := r.createResource(cr, m); err != nil {
		if errors.IsAlreadyExists(err) {
			e := &corev1.Service{ObjectMeta: m.ObjectMeta}
			if err := r.getResource(e); err != nil {
				return err
			}

			e.Spec.Ports = m.Spec.Ports
			e.SetLabels(m.GetLabels())

			if err := r.updateResource(cr, e); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (r *PrometheusAdapterReconciler) handleDeployment(cr *promv1.PrometheusAdapter) error {
	log := r.Log.WithValues("prometheusadapter", fmt.Sprintf("%s/%s", cr.GetNamespace(), cr.GetName()))
	f := NewFactory(cr)
	m, err := f.PrometheusAdapterDeployment()
	if err != nil {
		log.Error(err, "failed creating Deployment manifest")
		return err
	}

	if err = r.createResource(cr, m); err != nil {
		if errors.IsAlreadyExists(err) {
			e := &appsv1.Deployment{ObjectMeta: m.ObjectMeta}
			if err = r.getResource(e); err != nil {
				return err
			}

			e.SetLabels(m.GetLabels())
			e.Spec.Template = m.Spec.Template
			e.Spec.Replicas = m.Spec.Replicas

			if err = r.updateResource(cr, e); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func getTagFromImage(image string) string {
	partsOfImage := strings.Split(image, ":")
	return partsOfImage[len(partsOfImage)-1]
}

func getInstanceLabel(name, namespace string) string {
	label := fmt.Sprintf("%s-%s", name, namespace)
	if len(label) >= 63 {
		return strings.Trim(label[:63], "-")
	}
	return strings.Trim(label, "-")
}

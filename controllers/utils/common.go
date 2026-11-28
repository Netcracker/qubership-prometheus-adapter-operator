package utils

import (
    promv1 "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	corev1 "k8s.io/api/core/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const (
	PrometheusAdapterConfig = "prometheus-adapter-config"
)

func PrometheusAdapterIgnorePredicates() predicate.Predicate {
	return predicate.Funcs{
		DeleteFunc: func(e event.DeleteEvent) bool {
			return IsRelevantObject(e.Object) && !e.DeleteStateUnknown
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			switch obj := e.ObjectNew.(type) {
				case *corev1.ConfigMap:
					return obj.Name == PrometheusAdapterConfig
			  	case *promv1.PrometheusAdapter:
				  	// Ignore updates to CR status in which case metadata.Generation does not change
				  	return e.ObjectOld.GetGeneration() != obj.GetGeneration()
			  	default:
				  	return false
			}
		},
		CreateFunc: func(e event.CreateEvent) bool {
			return IsRelevantObject(e.Object)
		},
		GenericFunc: func(e event.GenericEvent) bool {
			return IsRelevantObject(e.Object)
		},
	}
}

func IsRelevantObject(obj client.Object) bool {
	switch obj := obj.(type) {
	case *corev1.ConfigMap:
		return obj.Name == PrometheusAdapterConfig
	case *promv1.PrometheusAdapter:
		return true
	default:
		return false
	}
}

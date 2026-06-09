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
	"fmt"
	"time"

	k8sapimeta "k8s.io/apimachinery/pkg/api/meta"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

var (
	// OneMinute represents one minute in milliseconds
	OneMinute = time.Minute

	// ThreeMinutes represents 3 minutes in milliseconds
	ThreeMinutes = 3 * time.Minute
)

// KindedNamespacedName holds K8s resource kind, namespace and name.
type KindedNamespacedName struct {
	Kind      string
	Namespace string
	Name      string
}

// GetKindedNamespacedName builds an instance KindedNamespacedName for given object.
func GetKindedNamespacedName(o k8sruntime.Object) *KindedNamespacedName {
	metaAccessor := k8sapimeta.NewAccessor()
	kind, _ := metaAccessor.Kind(o)
	namespace, _ := metaAccessor.Namespace(o)
	name, _ := metaAccessor.Name(o)

	return &KindedNamespacedName{
		Kind:      kind,
		Namespace: namespace,
		Name:      name,
	}
}

// NamespacedName builds string with object namespace and name separated with a slash.
func (m *KindedNamespacedName) NamespacedName() string {
	return fmt.Sprintf("%s/%s", m.Namespace, m.Name)
}

// K8sResource is an abstract k8s resource which can be reconciled
type K8sResource interface {
	k8sruntime.Object
	k8smetav1.Object

	GetObjectMeta() k8smetav1.Object
}

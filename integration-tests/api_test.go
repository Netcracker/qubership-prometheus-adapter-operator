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

package bdd_tests

import (
	"context"
	"fmt"
	"time"

	promv1 "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Custom resources", func() {
	ctx := context.Background()
	var namespace string

	BeforeEach(func() {
		namespace = fmt.Sprintf("prometheus-adapter-envtest-%d", time.Now().UnixNano())
		ns := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		err := k8sClient.Create(ctx, ns)
		if apierrors.IsAlreadyExists(err) {
			return
		}
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		err := k8sClient.Delete(ctx, &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		})
		if apierrors.IsNotFound(err) {
			return
		}
		Expect(err).NotTo(HaveOccurred())
	})

	It("stores and reads PrometheusAdapter resources", func() {
		replicas := int32(2)
		adapter := &promv1.PrometheusAdapter{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample",
				Namespace: namespace,
				Labels: map[string]string{
					"app.kubernetes.io/name": "prometheus-adapter",
				},
			},
			Spec: promv1.PrometheusAdapterSpec{
				Image:                 "example.com/prometheus-adapter:test",
				Replicas:              &replicas,
				PrometheusURL:         "https://prometheus.example.com",
				EnableCustomMetrics:   true,
				EnableResourceMetrics: true,
			},
		}

		Expect(k8sClient.Create(ctx, adapter)).To(Succeed())

		created := &promv1.PrometheusAdapter{}
		key := types.NamespacedName{Name: adapter.Name, Namespace: adapter.Namespace}
		Expect(k8sClient.Get(ctx, key, created)).To(Succeed())
		Expect(created.Spec.Image).To(Equal(adapter.Spec.Image))
		Expect(created.Spec.Replicas).To(HaveValue(Equal(replicas)))
		Expect(created.Spec.EnableCustomMetrics).To(BeTrue())
		Expect(created.Spec.EnableResourceMetrics).To(BeTrue())
	})

	It("stores and lists CustomScaleMetricRule resources", func() {
		rule := &promv1.CustomScaleMetricRule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "sample",
				Namespace: namespace,
			},
			Spec: promv1.CustomScaleMetricRuleSpec{
				Rules: []promv1.CustomMetricRuleConfig{
					{
						SeriesQuery: `http_requests_total{namespace!="",pod!=""}`,
						Resources: promv1.ResourceMapping{
							Overrides: map[string]promv1.GroupResource{
								"namespace": {Resource: "namespace"},
								"pod":       {Resource: "pod"},
							},
						},
						Name: promv1.NameMapping{
							Matches: "^(.*)_total$",
							As:      "${1}",
						},
						MetricsQuery: `sum(rate(<<.Series>>{<<.LabelMatchers>>}[2m])) by (<<.GroupBy>>)`,
					},
				},
			},
		}

		Expect(k8sClient.Create(ctx, rule)).To(Succeed())

		rules := &promv1.CustomScaleMetricRuleList{}
		Expect(k8sClient.List(ctx, rules, client.InNamespace(namespace))).To(Succeed())
		Expect(rules.Items).To(HaveLen(1))
		Expect(rules.Items[0].Spec.Rules).To(HaveLen(1))
		Expect(rules.Items[0].Spec.Rules[0].SeriesQuery).To(Equal(rule.Spec.Rules[0].SeriesQuery))
	})
})

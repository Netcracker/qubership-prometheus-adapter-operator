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
	api "github.com/Netcracker/qubership-prometheus-adapter-operator/api/v1"
	"github.com/Netcracker/qubership-prometheus-adapter-operator/controllers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manifests", func() {
	var (
		prometheusAdapterCr *api.PrometheusAdapter
		factory             *controllers.Factory
	)

	BeforeEach(func() {
		prometheusAdapterCr = &api.PrometheusAdapter{}
		factory = controllers.NewFactory(prometheusAdapterCr)
	})

	Describe("The prometheus-adapter Deployment manifest", func() {
		It("can be built", func() {
			d, err := factory.PrometheusAdapterDeployment()

			Expect(err).NotTo(HaveOccurred())
			Expect(d).ToNot(BeNil())
		})
	})

	Describe("The prometheus-adapter Service manifest", func() {
		It("can be built", func() {
			s, err := factory.PrometheusAdapterService()

			Expect(err).NotTo(HaveOccurred())
			Expect(s).ToNot(BeNil())
		})
	})

	Describe("The prometheus-adapter ServiceAccount manifest", func() {
		It("can be built", func() {
			sa, err := factory.PrometheusAdapterServiceAccount()

			Expect(err).NotTo(HaveOccurred())
			Expect(sa).ToNot(BeNil())
		})
	})

	Describe("The prometheus-adapter custom metrics ClusterRole manifest", func() {
		It("can be built", func() {
			clusterrole, err := factory.CustomMetricsClusterRole()

			Expect(err).NotTo(HaveOccurred())
			Expect(clusterrole).ToNot(BeNil())
		})
	})

	Describe("The prometheus-adapter custom metrics ClusterRoleBinding manifest", func() {
		It("can be built", func() {
			clusterrolebinding, err := factory.CustomMetricsClusterRoleBinding()

			Expect(err).NotTo(HaveOccurred())
			Expect(clusterrolebinding).ToNot(BeNil())
		})
	})
})

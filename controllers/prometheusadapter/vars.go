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
	pmodel "github.com/prometheus/common/model"
)

var (
	duration, _          = pmodel.ParseDuration("5m")
	defaultResourceRules = ResourceRules{
		CPU: ResourceRule{
			ContainerLabel: "container",
			ContainerQuery: "sum by (<<.GroupBy>>) (irate (container_cpu_usage_seconds_total{<<.LabelMatchers>>,container!=\"\",pod!=\"\"}[4m]))",
			NodeQuery:      "sum by (<<.GroupBy>>) (irate(node_cpu_seconds_total{<<.LabelMatchers>>}[4m]))",
			Resources: ResourceMapping{
				Overrides: map[string]GroupResource{
					"namespace": {
						Resource: "namespace",
					},
					"node": {
						Resource: "node",
					},
					"pod": {
						Resource: "pod",
					},
				},
			},
		},
		Memory: ResourceRule{
			ContainerLabel: "container",
			ContainerQuery: "sum by (<<.GroupBy>>) (container_memory_working_set_bytes{<<.LabelMatchers>>,container!=\"\",pod!=\"\"})",
			NodeQuery:      "sum by (<<.GroupBy>>) (node_memory_working_set_bytes{<<.LabelMatchers>>})",
			Resources: ResourceMapping{
				Overrides: map[string]GroupResource{
					"namespace": {
						Resource: "namespace",
					},
					"node": {
						Resource: "node",
					},
					"pod": {
						Resource: "pod",
					},
				},
			},
		},
		Window: duration,
	}
)

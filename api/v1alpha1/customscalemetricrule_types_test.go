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

package v1alpha1

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCustomScaleMetricRuleList_ItemsToString(t *testing.T) {
	tests := []struct {
		name     string
		items    []CustomScaleMetricRule
		expected string
	}{
		{
			name:     "empty list",
			items:    []CustomScaleMetricRule{},
			expected: "",
		},
		{
			name: "single item",
			items: []CustomScaleMetricRule{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule1",
						Namespace: "ns1",
					},
				},
			},
			expected: "ns1/rule1",
		},
		{
			name: "multiple items",
			items: []CustomScaleMetricRule{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule1",
						Namespace: "ns1",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule2",
						Namespace: "ns2",
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "rule3",
						Namespace: "ns1",
					},
				},
			},
			expected: "ns1/rule1, ns2/rule2, ns1/rule3",
		},
		{
			name: "item without namespace",
			items: []CustomScaleMetricRule{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "rule1",
					},
				},
			},
			expected: "/rule1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := &CustomScaleMetricRuleList{
				Items: tt.items,
			}
			result := list.ItemsToString()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

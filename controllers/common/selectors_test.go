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

package common

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMatchesSelector(t *testing.T) {
	tests := []struct {
		name     string
		object   metav1.Object
		selector *metav1.LabelSelector
		expected bool
		hasError bool
	}{
		{
			name: "empty selector matches any object",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test"},
			},
			selector: &metav1.LabelSelector{},
			expected: true,
			hasError: false,
		},
		{
			name: "matching labels",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test", "env": "prod"},
			},
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "test"},
			},
			expected: true,
			hasError: false,
		},
		{
			name: "non-matching labels",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test"},
			},
			selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "other"},
			},
			expected: false,
			hasError: false,
		},
		{
			name: "invalid selector",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test"},
			},
			selector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "app",
						Operator: "InvalidOperator",
						Values:   []string{"test"},
					},
				},
			},
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := matchesSelector(tt.object, tt.selector)
			if tt.hasError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMatchAll(t *testing.T) {
	tests := []struct {
		name      string
		object    metav1.Object
		selectors []*metav1.LabelSelector
		expected  bool
		hasError  bool
	}{
		{
			name: "empty selectors",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test"},
			},
			selectors: []*metav1.LabelSelector{},
			expected:  false,
			hasError:  false,
		},
		{
			name: "all selectors match",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test", "env": "prod"},
			},
			selectors: []*metav1.LabelSelector{
				{MatchLabels: map[string]string{"app": "test"}},
				{MatchLabels: map[string]string{"env": "prod"}},
			},
			expected: true,
			hasError: false,
		},
		{
			name: "one selector matches",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test"},
			},
			selectors: []*metav1.LabelSelector{
				{MatchLabels: map[string]string{"app": "test"}},
				{MatchLabels: map[string]string{"env": "prod"}},
			},
			expected: true,
			hasError: false,
		},
		{
			name: "invalid selector",
			object: &metav1.ObjectMeta{
				Labels: map[string]string{"app": "test"},
			},
			selectors: []*metav1.LabelSelector{
				{
					MatchExpressions: []metav1.LabelSelectorRequirement{
						{
							Key:      "app",
							Operator: "InvalidOperator",
							Values:   []string{"test"},
						},
					},
				},
			},
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MatchAll(tt.object, tt.selectors)
			if tt.hasError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

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

package controllers

import (
	"testing"
)

func TestKindedNamespacedName_NamespacedName(t *testing.T) {
	tests := []struct {
		name     string
		input    KindedNamespacedName
		expected string
	}{
		{
			name: "with namespace",
			input: KindedNamespacedName{
				Namespace: "test-ns",
				Name:      "test-name",
			},
			expected: "test-ns/test-name",
		},
		{
			name: "without namespace",
			input: KindedNamespacedName{
				Name: "test-name",
			},
			expected: "/test-name",
		},
		{
			name:     "empty",
			input:    KindedNamespacedName{},
			expected: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.NamespacedName()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

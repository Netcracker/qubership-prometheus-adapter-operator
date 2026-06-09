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

package common

import (
	apimachinerymetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerylabels "k8s.io/apimachinery/pkg/labels"
)

// matchesSelector checks if specified object matches specified label selector.
// Empty selector matches any object.
func matchesSelector(o apimachinerymetav1.Object, s *apimachinerymetav1.LabelSelector) (bool, error) {
	selector, err := apimachinerymetav1.LabelSelectorAsSelector(s)
	if err != nil {
		return false, err
	}

	return selector.Empty() || selector.Matches(apimachinerylabels.Set(o.GetLabels())), nil
}

// MatchAll checks if specified object matches specified label selector.
func MatchAll(o apimachinerymetav1.Object, s []*apimachinerymetav1.LabelSelector) (bool, error) {
	r := false

	for _, selector := range s {
		match, err := matchesSelector(o, selector)
		if err != nil {
			return false, err
		}
		r = r || match
	}

	return r, nil
}

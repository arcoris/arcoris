// Copyright 2026 The ARCORIS Authors.
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

package objectquery

import "testing"

func TestLabelSelectorEmptyMatchesEverything(t *testing.T) {
	var selector LabelSelector

	if !selector.IsZero() {
		t.Fatal("zero label selector is not zero")
	}
	if !selector.match(testItem("system", "worker", nil, nil)) {
		t.Fatal("zero label selector did not match")
	}
	requireNoError(t, selector.Validate())
}

func mustLabelSelector(t *testing.T, requirements ...LabelRequirement) LabelSelector {
	t.Helper()
	selector, err := NewLabelSelector(requirements...)
	requireNoError(t, err)

	return selector
}

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

func TestAnnotationSelectorEmptyMatchesEverything(t *testing.T) {
	var selector AnnotationSelector

	if !selector.IsZero() {
		t.Fatal("zero annotation selector is not zero")
	}
	if !selector.match(testItem("system", "worker", nil, nil)) {
		t.Fatal("zero annotation selector did not match")
	}
	requireNoError(t, selector.Validate())
}

func mustAnnotationSelector(t *testing.T, requirements ...AnnotationRequirement) AnnotationSelector {
	t.Helper()
	selector, err := NewAnnotationSelector(requirements...)
	requireNoError(t, err)

	return selector
}

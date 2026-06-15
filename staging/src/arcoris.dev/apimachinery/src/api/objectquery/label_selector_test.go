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
	if selector.Requirements() != nil {
		t.Fatalf("Requirements() = %#v; want nil", selector.Requirements())
	}
}

func TestLabelSelectorRequirementsAccessorCanonicalOrder(t *testing.T) {
	selector := mustLabelSelector(
		t,
		mustLabelEquals(t, "tier", "backend"),
		mustLabelIn(t, "env", "qa", "prod"),
	)

	requirements := selector.Requirements()
	if len(requirements) != 2 {
		t.Fatalf("len = %d; want 2", len(requirements))
	}
	if requirements[0].Key() != "env" || requirements[1].Key() != "tier" {
		t.Fatalf("requirement order = %q, %q; want env, tier", requirements[0].Key(), requirements[1].Key())
	}
}

func TestLabelSelectorRequirementsDefensiveCopy(t *testing.T) {
	selector := mustLabelSelector(t, mustLabelIn(t, "env", "qa", "prod"))
	requirements := selector.Requirements()
	requirements[0].req.values[0] = "mutated"
	requirements = append(requirements, mustLabelEquals(t, "tier", "backend"))

	next := selector.Requirements()
	if len(next) != 1 {
		t.Fatalf("len = %d; want 1", len(next))
	}
	requireStrings(t, next[0].Values(), "prod", "qa")
}

func mustLabelSelector(t *testing.T, requirements ...LabelRequirement) LabelSelector {
	t.Helper()
	selector, err := NewLabelSelector(requirements...)
	requireNoError(t, err)

	return selector
}

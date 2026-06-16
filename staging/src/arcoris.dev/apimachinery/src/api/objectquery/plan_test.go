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

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestPlanRangeConstraintsReturnsDetachedConstraints verifies callers cannot
// mutate predicate plan literals through iteration.
func TestPlanRangeConstraintsReturnsDetachedConstraints(t *testing.T) {
	predicate := mustPredicate(t, mustQ(LabelIn("env", "prod", "qa")))

	var first Constraint
	predicate.RangeConstraints(func(constraint Constraint) bool {
		first = constraint
		first.Values[0] = value.StringValue("mutated")
		return false
	})

	predicate.RangeConstraints(func(constraint Constraint) bool {
		got, _ := constraint.Values[0].AsString()
		if got != "prod" {
			t.Fatalf("constraint value = %q; want prod", got)
		}
		return false
	})
}

// TestPlanCloneDetachesLiteralSlices verifies the Plan accessor has the same
// defensive-copy behavior as RangeConstraints.
func TestPlanCloneDetachesLiteralSlices(t *testing.T) {
	predicate := mustPredicate(t, mustQ(LabelIn("env", "prod", "qa")))
	plan := predicate.Plan()

	plan.constraints[0].Values[0] = value.StringValue("mutated")

	gotPlan := predicate.Plan()
	got, _ := gotPlan.constraints[0].Values[0].AsString()
	if got != "prod" {
		t.Fatalf("cloned plan value = %q; want prod", got)
	}
}

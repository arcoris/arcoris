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

package objectindex

import "testing"

func TestCandidatePlanUnconstrainedIncludesEveryPosition(t *testing.T) {
	var plan candidatePlan

	for _, pos := range []int{-1, 0, 1, 100} {
		if !plan.includes(pos) {
			t.Fatalf("includes(%d) = false; want true for unconstrained plan", pos)
		}
	}
}

func TestCandidatePlanFirstConstraintDefinesCandidateUniverse(t *testing.T) {
	plan := candidatePlan{}.constrain(candidateSetFromPositions(4, []int{1, 3}))

	if !plan.constrained {
		t.Fatal("constrained = false; want true")
	}
	requirePlanIncludesOnly(t, plan, 4, 1, 3)
}

func TestCandidatePlanAdditionalConstraintsIntersect(t *testing.T) {
	plan := candidatePlan{}.
		constrain(candidateSetFromPositions(5, []int{0, 1, 3})).
		constrain(candidateSetFromPositions(5, []int{1, 2, 3}))

	requirePlanIncludesOnly(t, plan, 5, 1, 3)
}

func requirePlanIncludesOnly(t *testing.T, plan candidatePlan, size int, positions ...int) {
	t.Helper()

	want := make(map[int]struct{}, len(positions))
	for _, pos := range positions {
		want[pos] = struct{}{}
	}
	for pos := 0; pos < size; pos++ {
		_, present := want[pos]
		if got := plan.includes(pos); got != present {
			t.Fatalf("includes(%d) = %v; want %v", pos, got, present)
		}
	}
}

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

func TestCandidateSetFromPositions(t *testing.T) {
	set := candidateSetFromPositions(5, []int{-1, 0, 2, 2, 5})

	requireCandidateMembership(t, set, 5, 0, 2)
}

func TestUnionCandidateSets(t *testing.T) {
	set := unionCandidateSets(5, []int{0, 3}, []int{1, 3}, []int{-1, 5})

	requireCandidateMembership(t, set, 5, 0, 1, 3)
}

func TestIntersectCandidateSets(t *testing.T) {
	left := candidateSetFromPositions(5, []int{0, 1, 3})
	right := candidateSetFromPositions(5, []int{1, 2, 3})

	got := intersectCandidateSets(left, right)

	requireCandidateMembership(t, got, 5, 1, 3)
	requireCandidateMembership(t, left, 5, 0, 1, 3)
	requireCandidateMembership(t, right, 5, 1, 2, 3)
}

func TestCandidateSetHasRejectsOutOfRange(t *testing.T) {
	set := candidateSetFromPositions(2, []int{1})

	if set.has(-1) {
		t.Fatal("has(-1) = true; want false")
	}
	if set.has(2) {
		t.Fatal("has(2) = true; want false")
	}
}

func requireCandidateMembership(t *testing.T, set candidateSet, size int, positions ...int) {
	t.Helper()

	want := make(map[int]struct{}, len(positions))
	for _, pos := range positions {
		want[pos] = struct{}{}
	}
	for pos := 0; pos < size; pos++ {
		_, present := want[pos]
		if got := set.has(pos); got != present {
			t.Fatalf("has(%d) = %v; want %v", pos, got, present)
		}
	}
}

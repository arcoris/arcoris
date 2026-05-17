/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package runner

import "testing"

func TestCompactUsedPartialsRemovesInactiveSlots(t *testing.T) {
	partials := []string{"w0", "idle", "w2", "idle", "w4"}
	used := []bool{true, false, true, false, true}
	got := compactUsedPartials(partials, used)
	want := []string{"w0", "w2", "w4"}
	if len(got) != len(want) {
		t.Fatalf("compactUsedPartials length = %d, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("compactUsedPartials() = %#v, want %#v", got, want)
		}
	}
}

func TestCompactUsedPartialsPreservesWorkerOrder(t *testing.T) {
	partials := []int{10, 20, 30, 40}
	used := []bool{false, true, true, false}
	got := compactUsedPartials(partials, used)
	if len(got) != 2 || got[0] != 20 || got[1] != 30 {
		t.Fatalf("compactUsedPartials() = %#v, want [20 30]", got)
	}
}

func TestCompactUsedPartialsReturnsEmptyForNoActiveSlots(t *testing.T) {
	partials := []int{1, 2, 3}
	used := []bool{false, false, false}
	got := compactUsedPartials(partials, used)
	if len(got) != 0 {
		t.Fatalf("compactUsedPartials() length = %d, want 0", len(got))
	}
}

func TestCompactUsedPartialsPanicsOnLengthMismatch(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("compactUsedPartials did not panic on mismatched lengths")
		}
	}()
	_ = compactUsedPartials([]int{1, 2}, []bool{true})
}

func TestCompactUsedPartialsAllocatesNothing(t *testing.T) {
	partials := []int{1, 2, 3, 4}
	used := []bool{true, false, true, false}
	allocs := testing.AllocsPerRun(100, func() {
		_ = compactUsedPartials(partials, used)
	})
	if allocs != 0 {
		t.Fatalf("compactUsedPartials allocs = %f, want 0", allocs)
	}
}

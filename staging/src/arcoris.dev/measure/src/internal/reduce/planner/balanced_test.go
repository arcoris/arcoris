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


package planner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce/core"
)

func TestBalancedSmallInputFallsBackToSingleRange(t *testing.T) {
	got := Balanced(100, core.Options{Workers: 4, MinItemsPerWorker: 64}, nil)
	if len(got) != 1 || got[0].Start != 0 || got[0].End != 100 {
		t.Fatalf("Balanced fallback = %#v", got)
	}
}

func TestBalancedCoversInputWithoutGapsOrOverlap(t *testing.T) {
	got := Balanced(1000, core.Options{Workers: 4, MinItemsPerWorker: 100}, nil)
	assertPlanCovers(t, got, 1000)
	if len(got) != 4 {
		t.Fatalf("range count = %d, want 4", len(got))
	}
}

func TestBalancedRespectsMinItemsPerWorker(t *testing.T) {
	got := Balanced(1000, core.Options{Workers: 100, MinItemsPerWorker: 200}, nil)
	assertPlanCovers(t, got, 1000)
	if len(got) != 5 {
		t.Fatalf("range count = %d, want 5", len(got))
	}
}

func TestBalancedRespectsWorkers(t *testing.T) {
	got := Balanced(1000, core.Options{Workers: 3, MinItemsPerWorker: 1}, nil)
	assertPlanCovers(t, got, 1000)
	if len(got) != 3 {
		t.Fatalf("range count = %d, want 3", len(got))
	}
}

func TestBalancedRangeSizesDifferByAtMostOne(t *testing.T) {
	got := Balanced(1000, core.Options{Workers: 6, MinItemsPerWorker: 1}, nil)
	assertPlanCovers(t, got, 1000)
	minSize := got[0].Len()
	maxSize := got[0].Len()
	for _, r := range got[1:] {
		if r.Len() < minSize {
			minSize = r.Len()
		}
		if r.Len() > maxSize {
			maxSize = r.Len()
		}
	}
	if maxSize-minSize > 1 {
		t.Fatalf("range sizes differ by %d, want <= 1: %#v", maxSize-minSize, got)
	}
}

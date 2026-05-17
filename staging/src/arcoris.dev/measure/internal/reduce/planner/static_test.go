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

package planner

import (
	"testing"

	"arcoris.dev/measure/internal/reduce"
)

func TestStaticSmallInputFallsBackToSingleRange(t *testing.T) {
	got := Static(100, reduce.Options{Workers: 4, MinItemsPerWorker: 64}, nil)
	if len(got) != 1 || got[0].Start != 0 || got[0].End != 100 {
		t.Fatalf("Static fallback = %#v", got)
	}
}

func TestStaticCoversInputWithoutGapsOrOverlap(t *testing.T) {
	got := Static(1000, reduce.Options{Workers: 4, MinItemsPerWorker: 100}, nil)
	assertPlanCovers(t, got, 1000)
	if len(got) != 4 {
		t.Fatalf("range count = %d, want 4", len(got))
	}
}

func TestStaticRespectsMinItemsPerWorker(t *testing.T) {
	got := Static(1000, reduce.Options{Workers: 100, MinItemsPerWorker: 200}, nil)
	assertPlanCovers(t, got, 1000)
	if len(got) != 5 {
		t.Fatalf("range count = %d, want 5", len(got))
	}
}

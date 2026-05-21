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

// assertPlanCovers verifies that ranges exactly cover [0:n) without empty
// ranges, gaps, or overlaps.
func assertPlanCovers(
	t *testing.T,
	ranges []core.Range,
	n int,
) {
	t.Helper()
	pos := 0
	for i, r := range ranges {
		if r.Empty() {
			t.Fatalf("range[%d] is empty: %#v", i, r)
		}
		if r.Start != pos {
			t.Fatalf("range[%d].Start = %d, want %d; plan=%#v", i, r.Start, pos, ranges)
		}
		if r.End <= r.Start {
			t.Fatalf("range[%d] inverted: %#v", i, r)
		}
		pos = r.End
	}
	if pos != n {
		t.Fatalf("covered end = %d, want %d; plan=%#v", pos, n, ranges)
	}
}

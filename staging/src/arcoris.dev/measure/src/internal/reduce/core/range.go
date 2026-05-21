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


package core

// Range identifies a half-open index interval [Start, End).
//
// Range uses int indexes because reducers operate over Go slices and slice
// indexes are int. Planner implementations must emit valid, non-empty ranges in
// increasing order. Helper methods still treat inverted ranges as empty so tests
// and diagnostics can describe malformed plans without panicking.
type Range struct {
	// Start is the inclusive lower index of the interval.
	Start int

	// End is the exclusive upper index of the interval.
	End int
}

// Len returns the number of indexes covered by r.
func (r Range) Len() int {
	if r.End <= r.Start {
		return 0
	}
	return r.End - r.Start
}

// Empty reports whether r covers no indexes.
func (r Range) Empty() bool { return r.Len() == 0 }

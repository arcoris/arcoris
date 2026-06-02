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

package valuevalidation

import "testing"

func TestIntegerLimitsZeroValueHasNoBounds(t *testing.T) {
	var limits integerLimits[int64]

	if limits.lower.set {
		t.Fatalf("lower bound is set")
	}
	if limits.upper.set {
		t.Fatalf("upper bound is set")
	}
}

func TestIntegerBoundStoresValueAndPresence(t *testing.T) {
	bound := integerBound[uint64]{
		value: 42,
		set:   true,
	}

	if !bound.set {
		t.Fatalf("set = false")
	}
	if got := bound.value; got != 42 {
		t.Fatalf("value = %d, want 42", got)
	}
}

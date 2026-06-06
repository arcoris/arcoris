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

package capacity_test

import (
	"math"
	"testing"

	"arcoris.dev/capacity"
)

func TestVectorCheckedAddAndSub(t *testing.T) {
	left := vector(t, entry("memory_bytes", 8), entry("worker_slots", 2))
	right := vector(t, entry("queue_slots", 1), entry("worker_slots", 3))

	sum, ok := left.CheckedAdd(right)
	if !ok {
		t.Fatal("CheckedAdd() returned ok=false")
	}
	requireVector(t, sum, entry("memory_bytes", 8), entry("queue_slots", 1), entry("worker_slots", 5))

	diff, ok := sum.CheckedSub(vector(t, entry("memory_bytes", 8), entry("worker_slots", 5)))
	if !ok {
		t.Fatal("CheckedSub() returned ok=false")
	}
	requireVector(t, diff, entry("queue_slots", 1))
}

func TestVectorCheckedAddAndSubFailures(t *testing.T) {
	overflow := vector(t, capacity.Entry{
		Resource: capacity.MustResource("worker_slots"),
		Amount:   capacity.Amount(math.MaxUint64),
	})
	if _, ok := overflow.CheckedAdd(vector(t, entry("worker_slots", 1))); ok {
		t.Fatal("CheckedAdd() overflow returned ok=true")
	}

	diff := vector(t, entry("queue_slots", 1))
	if _, ok := diff.CheckedSub(vector(t, entry("worker_slots", 1))); ok {
		t.Fatal("CheckedSub() missing resource returned ok=true")
	}
}

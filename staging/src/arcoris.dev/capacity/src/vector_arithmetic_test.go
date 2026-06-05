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
)

func TestVectorCheckedAdd(t *testing.T) {
	t.Parallel()

	left := vector(t, entry("memory_bytes", 8), entry("worker_slots", 2))
	right := vector(t, entry("queue_slots", 3), entry("worker_slots", 4))

	got, ok := left.CheckedAdd(right)
	if !ok {
		t.Fatal("CheckedAdd() returned ok=false")
	}
	requireVector(t, got, entry("memory_bytes", 8), entry("queue_slots", 3), entry("worker_slots", 6))
	requireVector(t, left, entry("memory_bytes", 8), entry("worker_slots", 2))
}

func TestVectorCheckedAddOverflow(t *testing.T) {
	t.Parallel()

	left := vector(t, entry("worker_slots", math.MaxUint64))
	right := vector(t, entry("worker_slots", 1))

	if _, ok := left.CheckedAdd(right); ok {
		t.Fatal("CheckedAdd() overflow returned ok=true")
	}
}

func TestVectorCheckedSub(t *testing.T) {
	t.Parallel()

	left := vector(t, entry("memory_bytes", 8), entry("worker_slots", 6))
	right := vector(t, entry("memory_bytes", 8), entry("worker_slots", 2))

	got, ok := left.CheckedSub(right)
	if !ok {
		t.Fatal("CheckedSub() returned ok=false")
	}
	requireVector(t, got, entry("worker_slots", 4))
	requireVector(t, left, entry("memory_bytes", 8), entry("worker_slots", 6))
}

func TestVectorCheckedSubUnderflow(t *testing.T) {
	t.Parallel()

	left := vector(t, entry("worker_slots", 1))
	right := vector(t, entry("memory_bytes", 1))

	if _, ok := left.CheckedSub(right); ok {
		t.Fatal("CheckedSub() missing resource returned ok=true")
	}
}

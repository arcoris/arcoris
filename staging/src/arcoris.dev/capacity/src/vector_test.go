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
	"testing"

	"arcoris.dev/capacity"
)

func TestVectorAccessAndCopySafety(t *testing.T) {
	v := vector(t, entry("worker_slots", 2), entry("memory_bytes", 8))
	requireVector(t, v, entry("memory_bytes", 8), entry("worker_slots", 2))

	entries := v.Entries()
	entries[0] = entry("queue_slots", 1)
	requireVector(t, v, entry("memory_bytes", 8), entry("worker_slots", 2))

	if !v.Has(capacity.MustResource("worker_slots")) {
		t.Fatal("Has(worker_slots) = false")
	}
	if got := v.Amount(capacity.MustResource("worker_slots")); got != 2 {
		t.Fatalf("Amount(worker_slots) = %d, want 2", got)
	}
	if v.Has(capacity.MustResource("queue_slots")) {
		t.Fatal("Has(queue_slots) = true")
	}
	if got := v.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}

	other := vector(t, entry("worker_slots", 2), entry("memory_bytes", 8))
	if !v.Equal(other) {
		t.Fatal("Equal() = false for matching canonical vectors")
	}
}

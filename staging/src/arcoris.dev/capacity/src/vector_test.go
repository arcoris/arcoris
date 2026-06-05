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

func TestVectorEntriesAreFreshCopies(t *testing.T) {
	t.Parallel()

	v := vector(t, entry("worker_slots", 2))
	entries := v.Entries()
	entries[0] = entry("worker_slots", 99)

	if got := v.Amount(capacity.MustResource("worker_slots")); got != 2 {
		t.Fatalf("Amount() = %d, want 2", got)
	}
}

func TestVectorLookupAndEquality(t *testing.T) {
	t.Parallel()

	v := vector(t, entry("worker_slots", 2), entry("memory_bytes", 8))
	workerSlots := capacity.MustResource("worker_slots")
	queueSlots := capacity.MustResource("queue_slots")

	if !v.Has(workerSlots) || v.Has(queueSlots) {
		t.Fatalf("Has() mismatch")
	}
	if got := v.Amount(workerSlots); got != 2 {
		t.Fatalf("Amount(worker_slots) = %d, want 2", got)
	}
	if got := v.Amount(queueSlots); got != 0 {
		t.Fatalf("Amount(queue_slots) = %d, want 0", got)
	}
	if !v.Equal(vector(t, entry("memory_bytes", 8), entry("worker_slots", 2))) {
		t.Fatal("canonical vectors were not equal")
	}
}

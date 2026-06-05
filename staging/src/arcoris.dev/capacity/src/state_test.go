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

func TestStateReserveSuccessIsAllOrNothing(t *testing.T) {
	t.Parallel()

	state := capacity.MustState(
		vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)),
		vector(t, entry("worker_slots", 1)),
	)
	next, result := state.Reserve(demand(t, entry("memory_bytes", 3), entry("worker_slots", 2)))
	if !result.Reserved() {
		t.Fatalf("Reserve() status = %s, want reserved", result.Status)
	}
	requireVector(t, next.Reserved, entry("memory_bytes", 3), entry("worker_slots", 3))
	requireVector(t, state.Reserved, entry("worker_slots", 1))
}

func TestStateReserveRefusalDoesNotMutate(t *testing.T) {
	t.Parallel()

	state := capacity.MustState(
		vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)),
		vector(t, entry("worker_slots", 1)),
	)
	next, result := state.Reserve(demand(t, entry("memory_bytes", 3), entry("worker_slots", 4)))
	if result.Status != capacity.ReserveStatusInsufficient {
		t.Fatalf("Reserve() status = %s, want insufficient", result.Status)
	}
	if !next.Reserved.Equal(state.Reserved) {
		t.Fatalf("state mutated on refusal: %#v", next.Reserved.Entries())
	}
}

func TestStateReserveDebtDoesNotBlockUnrelatedResources(t *testing.T) {
	t.Parallel()

	state := capacity.MustState(
		vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)),
		vector(t, entry("memory_bytes", 10)),
	)
	next, result := state.Reserve(demand(t, entry("worker_slots", 2)))
	if !result.Reserved() {
		t.Fatalf("Reserve(worker_slots) status = %s, want reserved", result.Status)
	}
	requireVector(t, next.Reserved, entry("memory_bytes", 10), entry("worker_slots", 2))

	_, result = state.Reserve(demand(t, entry("memory_bytes", 1)))
	if result.Status != capacity.ReserveStatusDebt {
		t.Fatalf("Reserve(memory_bytes) status = %s, want debt", result.Status)
	}
}

func TestStateRelease(t *testing.T) {
	t.Parallel()

	state := capacity.MustState(
		vector(t, entry("worker_slots", 4)),
		vector(t, entry("worker_slots", 3)),
	)
	next, ok := state.Release(demand(t, entry("worker_slots", 2)))
	if !ok {
		t.Fatal("Release() returned ok=false")
	}
	requireVector(t, next.Reserved, entry("worker_slots", 1))

	unchanged, ok := state.Release(demand(t, entry("worker_slots", 4)))
	if ok {
		t.Fatal("Release() underflow returned ok=true")
	}
	if !unchanged.Reserved.Equal(state.Reserved) {
		t.Fatal("Release() underflow mutated state")
	}
}

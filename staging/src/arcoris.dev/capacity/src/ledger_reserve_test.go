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
	"sync"
	"testing"

	"arcoris.dev/capacity"
)

func TestLedgerReserveSuccess(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)))
	before := ledger.Revision()
	result := ledger.TryReserve(demand(t, entry("memory_bytes", 3), entry("worker_slots", 2)))
	if !result.Reserved() {
		t.Fatalf("TryReserve() status = %s, want reserved", result.Status)
	}
	if result.Reservation == nil {
		t.Fatal("successful reservation is nil")
	}
	if result.Snapshot.Revision == before {
		t.Fatal("successful reserve did not advance revision")
	}
	requireVector(t, result.Snapshot.Value.Reserved, entry("memory_bytes", 3), entry("worker_slots", 2))
}

func TestLedgerReserveRefusalDoesNotMutateOrAdvance(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)))
	first := ledger.TryReserve(demand(t, entry("memory_bytes", 3)))
	if !first.Reserved() {
		t.Fatalf("first reserve status = %s", first.Status)
	}
	before := ledger.Snapshot()

	denied := ledger.TryReserve(demand(t, entry("memory_bytes", 1), entry("worker_slots", 5)))
	if denied.Status != capacity.ReserveStatusInsufficient {
		t.Fatalf("denied status = %s, want insufficient", denied.Status)
	}
	if denied.Reservation != nil {
		t.Fatal("denied reservation is non-nil")
	}
	if denied.Snapshot.Revision != before.Revision {
		t.Fatal("failed reserve advanced revision")
	}
	requireVector(t, denied.Snapshot.Value.Reserved, entry("memory_bytes", 3))
}

func TestLedgerUnknownResourceDoesNotPartiallyReserve(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 4)))
	denied := ledger.TryReserve(demand(t, entry("queue_slots", 1), entry("worker_slots", 2)))
	if denied.Status != capacity.ReserveStatusUnknownResource {
		t.Fatalf("status = %s, want unknown_resource", denied.Status)
	}
	requireVector(t, ledger.Snapshot().Value.Reserved)
}

func TestLedgerDebtBlocksDemandedResourceOnly(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)))
	memory := ledger.TryReserve(demand(t, entry("memory_bytes", 8)))
	if !memory.Reserved() {
		t.Fatalf("memory reserve status = %s", memory.Status)
	}
	overcommitted := ledger.SetLimits(vector(t, entry("memory_bytes", 4), entry("worker_slots", 4)))
	requireVector(t, overcommitted.Value.Debt, entry("memory_bytes", 4))

	worker := ledger.TryReserve(demand(t, entry("worker_slots", 2)))
	if !worker.Reserved() {
		t.Fatalf("worker reserve status = %s, want reserved", worker.Status)
	}

	denied := ledger.TryReserve(demand(t, entry("memory_bytes", 1)))
	if denied.Status != capacity.ReserveStatusDebt {
		t.Fatalf("memory reserve status = %s, want debt", denied.Status)
	}
}

func TestLedgerReserveMatchesScalarForOneResource(t *testing.T) {
	t.Parallel()

	resource := capacity.MustResource("worker_slots")
	multi := capacity.NewLedger(vector(t, entry("worker_slots", 4)))
	scalar := capacity.NewScalarLedger(4)

	multiResult := multi.TryReserve(demand(t, entry("worker_slots", 3)))
	scalarResult := scalar.TryReserve(3)
	if multiResult.Status != scalarResult.Status {
		t.Fatalf("status multi=%s scalar=%s", multiResult.Status, scalarResult.Status)
	}
	if got := multiResult.Snapshot.Value.Reserved.Amount(resource); got != scalarResult.Snapshot.Value.Reserved {
		t.Fatalf("reserved multi=%d scalar=%d", got, scalarResult.Snapshot.Value.Reserved)
	}

	_ = multi.SetLimits(vector(t, entry("worker_slots", 2)))
	_ = scalar.SetLimit(2)
	if got := multi.Snapshot().Value.Debt.Amount(resource); got != scalar.Snapshot().Value.Debt {
		t.Fatalf("debt multi=%d scalar=%d", got, scalar.Snapshot().Value.Debt)
	}
}

func TestLedgerConcurrentReserveDoesNotOverspend(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 32)))
	unit := demand(t, entry("worker_slots", 1))
	var wg sync.WaitGroup
	for i := 0; i < 128; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ledger.TryReserve(unit)
		}()
	}
	wg.Wait()

	snap := ledger.Snapshot()
	if got := snap.Value.Reserved.Amount(capacity.MustResource("worker_slots")); got != 32 {
		t.Fatalf("reserved = %d, want 32", got)
	}
}

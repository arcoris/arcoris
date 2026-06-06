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

func TestVectorLedgerTryReserveObserved(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t,
		entry("memory_bytes", 8),
		entry("worker_slots", 4),
	))

	reservation, observation, ok := ledger.TryReserveObserved(demand(t,
		entry("memory_bytes", 3),
		entry("worker_slots", 2),
	))
	if !ok {
		t.Fatalf("TryReserveObserved() refused: %#v", observation)
	}
	if reservation == nil {
		t.Fatal("TryReserveObserved() returned nil reservation")
	}
	if observation.Refusal != capacity.RefusalNone {
		t.Fatalf("Refusal = %s, want none", observation.Refusal)
	}
}

func TestVectorLedgerTryReserveRefusalDoesNotMutate(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t,
		entry("memory_bytes", 8),
		entry("worker_slots", 4),
	))

	reservation, ok := ledger.TryReserve(demand(t,
		entry("memory_bytes", 3),
		entry("worker_slots", 2),
	))
	if !ok {
		t.Fatal("initial reserve failed")
	}
	defer reservation.Release()

	before := ledger.Snapshot()
	_, refused, ok := ledger.TryReserveObserved(demand(t,
		entry("memory_bytes", 1),
		entry("worker_slots", 3),
	))
	if ok || refused.Refusal != capacity.RefusalInsufficient {
		t.Fatalf("refused = %#v ok=%v", refused, ok)
	}

	after := ledger.Snapshot()
	if !after.Value.Reserved.Equal(before.Value.Reserved) {
		t.Fatalf("reserved after refusal = %#v, want %#v", after.Value.Reserved, before.Value.Reserved)
	}
	if after.Revision != before.Revision {
		t.Fatalf("revision after refusal = %d, want %d", after.Revision, before.Revision)
	}
}

func TestVectorLedgerRawReserveAndReleaseRevision(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t,
		entry("worker_slots", 2),
	))
	unit := demand(t, entry("worker_slots", 1))
	initial := ledger.Revision()

	reservation, ok := ledger.TryReserve(unit)
	if !ok {
		t.Fatal("TryReserve() failed")
	}
	afterReserve := ledger.Revision()
	requireRevisionAdvanced(t, initial, afterReserve)

	extra, ok := ledger.TryReserve(demand(t, entry("worker_slots", 2)))
	if ok || extra != nil {
		t.Fatalf("extra TryReserve() = %v, %v; want nil, false", extra, ok)
	}
	requireRevisionEqual(t, afterReserve, ledger.Revision(), "failed vector TryReserve")

	if !reservation.TryRelease() {
		t.Fatal("VectorReservation.TryRelease() failed")
	}
	afterRelease := ledger.Revision()
	requireRevisionAdvanced(t, afterReserve, afterRelease)

	if reservation.TryRelease() {
		t.Fatal("second VectorReservation.TryRelease() succeeded")
	}
	requireRevisionEqual(t, afterRelease, ledger.Revision(), "second vector TryRelease")
}

func TestVectorLedgerDebtOnlyBlocksAffectedResources(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t,
		entry("memory_bytes", 8),
		entry("worker_slots", 4),
	))

	memory, ok := ledger.TryReserve(demand(t, entry("memory_bytes", 8)))
	if !ok {
		t.Fatal("memory reserve failed")
	}

	ledger.SetLimits(vector(t,
		entry("memory_bytes", 4),
		entry("worker_slots", 4),
	))

	worker, ok := ledger.TryReserve(demand(t, entry("worker_slots", 2)))
	if !ok {
		t.Fatal("worker reserve was blocked by unrelated memory debt")
	}

	_, debt, ok := ledger.TryReserveObserved(demand(t, entry("memory_bytes", 1)))
	if ok {
		t.Fatalf("memory reserve succeeded despite debt: %#v", debt)
	}
	if debt.Refusal != capacity.RefusalDebt {
		t.Fatalf("Refusal = %s, want debt", debt.Refusal)
	}

	worker.Release()
	memory.Release()
}

func TestVectorLedgerConcurrentReserveRelease(t *testing.T) {
	const limit = 32

	ledger := capacity.NewVectorLedger(vector(t,
		entry("worker_slots", limit),
	))
	unit := demand(t, entry("worker_slots", 1))
	var reservations sync.Map
	var wg sync.WaitGroup

	for i := 0; i < limit*4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			if reservation, ok := ledger.TryReserve(unit); ok {
				reservations.Store(i, reservation)
			}
		}(i)
	}
	wg.Wait()

	snapshot := ledger.Snapshot()
	if got := snapshot.Value.Reserved.Amount(capacity.MustResource("worker_slots")); got != limit {
		t.Fatalf("reserved = %d, want %d", got, limit)
	}

	reservations.Range(func(_, value any) bool {
		wg.Add(1)
		go func(reservation *capacity.VectorReservation) {
			defer wg.Done()
			_ = reservation.TryRelease()
		}(value.(*capacity.VectorReservation))

		return true
	})
	wg.Wait()

	snapshot = ledger.Snapshot()
	if !snapshot.Value.Reserved.IsZero() {
		t.Fatalf("reserved after releases = %#v", snapshot.Value.Reserved)
	}
}

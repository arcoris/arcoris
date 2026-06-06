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
	"arcoris.dev/snapshot"
)

func TestLedgerImplementsSnapshotSources(t *testing.T) {
	var _ snapshot.Source[capacity.Snapshot] = (*capacity.Ledger)(nil)
	var _ snapshot.RevisionSource = (*capacity.Ledger)(nil)
}

func TestLedgerSnapshotRevisionAndSetLimit(t *testing.T) {
	ledger := capacity.NewLedger(4)
	initial := ledger.Snapshot()

	if initial.Value != capacity.NewSnapshot(4, 0) {
		t.Fatalf("initial snapshot = %#v", initial.Value)
	}

	ledger.SetLimit(4)
	noChange := ledger.Snapshot()

	if noChange.Revision != initial.Revision {
		t.Fatalf("same SetLimit advanced revision: %d -> %d", initial.Revision, noChange.Revision)
	}

	changed := ledger.SetLimitObserved(2)

	if changed.Value != capacity.NewSnapshot(2, 0) {
		t.Fatalf("changed snapshot = %#v", changed.Value)
	}
	if changed.Revision == initial.Revision {
		t.Fatal("changed SetLimit did not advance revision")
	}
}

func TestLedgerRevisionAdvancesOnlyForSuccessfulMutations(t *testing.T) {
	ledger := capacity.NewLedger(2)
	initial := ledger.Revision()

	if !ledger.TryReserve(1) {
		t.Fatal("TryReserve(1) failed")
	}
	afterReserve := ledger.Revision()
	requireRevisionAdvanced(t, initial, afterReserve)

	if ledger.TryReserve(3) {
		t.Fatal("TryReserve(3) succeeded")
	}
	requireRevisionEqual(t, afterReserve, ledger.Revision(), "failed TryReserve")

	if !ledger.TryRelease(1) {
		t.Fatal("TryRelease(1) failed")
	}
	afterRelease := ledger.Revision()
	requireRevisionAdvanced(t, afterReserve, afterRelease)

	if ledger.TryRelease(1) {
		t.Fatal("TryRelease(1) succeeded without reservation")
	}
	requireRevisionEqual(t, afterRelease, ledger.Revision(), "failed TryRelease")

	reservation, ok := ledger.TryAcquire(1)
	if !ok {
		t.Fatal("TryAcquire(1) failed")
	}
	afterAcquire := ledger.Revision()
	requireRevisionAdvanced(t, afterRelease, afterAcquire)

	if !reservation.TryRelease() {
		t.Fatal("Reservation.TryRelease() failed")
	}
	afterReservationRelease := ledger.Revision()
	requireRevisionAdvanced(t, afterAcquire, afterReservationRelease)

	if reservation.TryRelease() {
		t.Fatal("second Reservation.TryRelease() succeeded")
	}
	requireRevisionEqual(t, afterReservationRelease, ledger.Revision(), "second reservation TryRelease")

	ledger.SetLimit(2)
	requireRevisionEqual(t, afterReservationRelease, ledger.Revision(), "same SetLimit")

	ledger.SetLimit(1)
	requireRevisionAdvanced(t, afterReservationRelease, ledger.Revision())
}

func TestLedgerConcurrentSetLimitReserveReleaseAndSnapshot(t *testing.T) {
	ledger := capacity.NewLedger(8)
	var wg sync.WaitGroup
	errs := make(chan string, 256)

	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 512; i++ {
				ledger.SetLimit(capacity.Amount((i + worker) % 16))
				if ledger.TryReserve(1) {
					ledger.Release(1)
				}

				snap := ledger.Snapshot()
				if !snap.Value.IsValid() {
					errs <- "invalid snapshot"
					return
				}
			}
		}(worker)
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatal(err)
	}
}

func TestLedgerLimitShrinkCreatesDebtUntilRelease(t *testing.T) {
	ledger := capacity.NewLedger(2)
	if !ledger.TryReserve(2) {
		t.Fatal("TryReserve(2) failed")
	}

	ledger.SetLimit(1)
	debt := ledger.Snapshot().Value
	if debt.Debt != 1 || debt.Reserved != 2 {
		t.Fatalf("snapshot after shrink = %#v", debt)
	}

	ledger.Release(1)
	cleared := ledger.Snapshot().Value
	if cleared.Debt != 0 || cleared.Reserved != 1 {
		t.Fatalf("snapshot after partial release = %#v", cleared)
	}

	ledger.Release(1)
	available := ledger.Snapshot().Value
	if available != capacity.NewSnapshot(1, 0) {
		t.Fatalf("snapshot after final release = %#v", available)
	}
}

func requireRevisionAdvanced(t *testing.T, before snapshot.Revision, after snapshot.Revision) {
	t.Helper()

	if after == before {
		t.Fatalf("revision did not advance: %d", before)
	}
}

func requireRevisionEqual(t *testing.T, want snapshot.Revision, got snapshot.Revision, operation string) {
	t.Helper()

	if got != want {
		t.Fatalf("%s advanced revision: got %d, want %d", operation, got, want)
	}
}

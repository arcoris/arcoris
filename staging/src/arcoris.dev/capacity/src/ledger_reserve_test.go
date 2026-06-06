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
	"sync/atomic"
	"testing"

	"arcoris.dev/capacity"
)

func TestLedgerTryReserveUsesRawAccounting(t *testing.T) {
	ledger := capacity.NewLedger(4)

	if !ledger.TryReserve(3) {
		t.Fatal("TryReserve(3) returned false")
	}
	if got := ledger.Snapshot().Value; got != capacity.NewSnapshot(4, 3) {
		t.Fatalf("snapshot after raw reserve = %#v", got)
	}

	ledger.Release(3)
	if got := ledger.Snapshot().Value; got != capacity.NewSnapshot(4, 0) {
		t.Fatalf("snapshot after raw release = %#v", got)
	}
}

func TestLedgerTryReserveObservedReportsDiagnostics(t *testing.T) {
	ledger := capacity.NewLedger(4)

	observation, ok := ledger.TryReserveObserved(3)
	if !ok {
		t.Fatal("TryReserveObserved(3) returned false")
	}
	if observation.Refusal != capacity.RefusalNone {
		t.Fatalf("Refusal = %s, want none", observation.Refusal)
	}
	if observation.Snapshot.Value != capacity.NewSnapshot(4, 3) {
		t.Fatalf("snapshot = %#v", observation.Snapshot.Value)
	}
}

func TestLedgerTryAcquireReturnsOwnedReservation(t *testing.T) {
	ledger := capacity.NewLedger(4)

	reservation, ok := ledger.TryAcquire(3)
	if !ok || reservation == nil {
		t.Fatalf("TryAcquire(3) = %v, %v", reservation, ok)
	}
	if got := reservation.Amount(); got != 3 {
		t.Fatalf("reservation amount = %d, want 3", got)
	}

	reservation.Release()
	if got := ledger.Snapshot().Value; got != capacity.NewSnapshot(4, 0) {
		t.Fatalf("snapshot after owned release = %#v", got)
	}
}

func TestLedgerTryAcquireObservedReportsDiagnostics(t *testing.T) {
	ledger := capacity.NewLedger(4)

	reservation, observation, ok := ledger.TryAcquireObserved(2)
	if !ok || reservation == nil {
		t.Fatalf("TryAcquireObserved(2) = %v, %#v, %v", reservation, observation, ok)
	}
	defer reservation.Release()

	if observation.Refusal != capacity.RefusalNone {
		t.Fatalf("Refusal = %s, want none", observation.Refusal)
	}
	if observation.Snapshot.Value != capacity.NewSnapshot(4, 2) {
		t.Fatalf("snapshot = %#v", observation.Snapshot.Value)
	}
}

func TestLedgerTryAcquireFailureReturnsNilReservation(t *testing.T) {
	ledger := capacity.NewLedger(1)

	reservation, ok := ledger.TryAcquire(2)
	if ok || reservation != nil {
		t.Fatalf("TryAcquire(2) = %v, %v; want nil, false", reservation, ok)
	}
	if got := ledger.Snapshot().Value; got.Reserved != 0 {
		t.Fatalf("reserved after failed acquire = %d, want 0", got.Reserved)
	}
}

func TestLedgerRefusalsDoNotMutate(t *testing.T) {
	ledger := capacity.NewLedger(2)

	if ledger.TryReserve(3) {
		t.Fatal("TryReserve(3) succeeded, want refusal")
	}
	if got := ledger.Snapshot().Value; got.Reserved != 0 {
		t.Fatalf("reserved after refusal = %d, want 0", got.Reserved)
	}

	insufficient, ok := ledger.TryReserveObserved(3)
	if ok || insufficient.Refusal != capacity.RefusalInsufficient {
		t.Fatalf("insufficient observation = %#v ok=%v", insufficient, ok)
	}

	if !ledger.TryReserve(2) {
		t.Fatal("TryReserve(2) failed")
	}
	ledger.SetLimit(1)

	debt, ok := ledger.TryReserveObserved(1)
	if ok || debt.Refusal != capacity.RefusalDebt || debt.Snapshot.Value.Debt != 1 {
		t.Fatalf("debt observation = %#v ok=%v", debt, ok)
	}

	ledger.Release(2)
}

func TestLedgerRawReleaseObservedReturnsSnapshot(t *testing.T) {
	ledger := capacity.NewLedger(2)
	if !ledger.TryReserve(2) {
		t.Fatal("TryReserve(2) failed")
	}

	snap := ledger.ReleaseObserved(1)
	if snap.Value != capacity.NewSnapshot(2, 1) {
		t.Fatalf("release snapshot = %#v", snap.Value)
	}

	ledger.Release(1)
}

func TestLedgerTryReleaseObservedReturnsOutcome(t *testing.T) {
	ledger := capacity.NewLedger(1)

	failed, ok := ledger.TryReleaseObserved(1)
	if ok || failed.Value.Reserved != 0 {
		t.Fatalf("failed TryReleaseObserved() = %#v, %v", failed.Value, ok)
	}

	if !ledger.TryReserve(1) {
		t.Fatal("TryReserve(1) failed")
	}
	released, ok := ledger.TryReleaseObserved(1)
	if !ok || released.Value.Reserved != 0 {
		t.Fatalf("successful TryReleaseObserved() = %#v, %v", released.Value, ok)
	}
}

func TestLedgerRawReleasePanicsOnUnderflow(t *testing.T) {
	ledger := capacity.NewLedger(1)

	requirePanicIs(t, capacity.ErrReservedUnderflow, func() {
		ledger.Release(1)
	})
}

func TestLedgerTryReleaseReportsUnderflow(t *testing.T) {
	ledger := capacity.NewLedger(1)

	if ledger.TryRelease(1) {
		t.Fatal("TryRelease(1) returned true, want false")
	}
	if got := ledger.Snapshot().Value; got.Reserved != 0 {
		t.Fatalf("reserved after failed TryRelease = %d, want 0", got.Reserved)
	}
}

func TestLedgerReserveAndReleasePanicsOnZeroAmount(t *testing.T) {
	ledger := capacity.NewLedger(1)

	requirePanicIs(t, capacity.ErrZeroAmount, func() { _ = ledger.TryReserve(0) })
	requirePanicIs(t, capacity.ErrZeroAmount, func() { ledger.Release(0) })
	requirePanicIs(t, capacity.ErrZeroAmount, func() { _, _ = ledger.TryAcquire(0) })
}

func TestLedgerConcurrentRawReserveDoesNotOverspend(t *testing.T) {
	const limit = 64

	ledger := capacity.NewLedger(limit)
	var successes atomic.Uint64
	var wg sync.WaitGroup

	for i := 0; i < limit*4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ledger.TryReserve(1) {
				successes.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := successes.Load(); got != limit {
		t.Fatalf("successful reserves = %d, want %d", got, limit)
	}
	if got := ledger.Snapshot().Value.Reserved; got != limit {
		t.Fatalf("reserved = %d, want %d", got, limit)
	}

	for i := 0; i < limit; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ledger.Release(1)
		}()
	}
	wg.Wait()

	if got := ledger.Snapshot().Value; got.Reserved != 0 || !got.IsValid() {
		t.Fatalf("snapshot after releases = %#v", got)
	}
}

func TestLedgerConcurrentTryReleaseSucceedsOnce(t *testing.T) {
	ledger := capacity.NewLedger(1)
	if !ledger.TryReserve(1) {
		t.Fatal("TryReserve(1) failed")
	}

	var successes atomic.Uint64
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ledger.TryRelease(1) {
				successes.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := successes.Load(); got != 1 {
		t.Fatalf("successful releases = %d, want 1", got)
	}
	if got := ledger.Snapshot().Value.Reserved; got != 0 {
		t.Fatalf("reserved = %d, want 0", got)
	}
}

func TestLedgerConcurrentTryReleaseDoesNotUnderflowOrAdvanceOnFailures(t *testing.T) {
	const reserved = 32

	ledger := capacity.NewLedger(reserved)
	if !ledger.TryReserve(reserved) {
		t.Fatal("TryReserve(reserved) failed")
	}
	before := ledger.Revision()

	var successes atomic.Uint64
	var wg sync.WaitGroup
	for i := 0; i < reserved*4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if ledger.TryRelease(1) {
				successes.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := successes.Load(); got != reserved {
		t.Fatalf("successful releases = %d, want %d", got, reserved)
	}
	after := ledger.Revision()
	if got := uint64(after - before); got != reserved {
		t.Fatalf("revision advances = %d, want %d", got, reserved)
	}
	if got := ledger.Snapshot().Value.Reserved; got != 0 {
		t.Fatalf("reserved = %d, want 0", got)
	}

	if ledger.TryRelease(1) {
		t.Fatal("extra TryRelease(1) succeeded")
	}
	requireRevisionEqual(t, after, ledger.Revision(), "extra TryRelease")
}

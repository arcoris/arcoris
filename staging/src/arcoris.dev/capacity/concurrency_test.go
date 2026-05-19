/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package capacity_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/capacity"
)

func TestLedgerConcurrentReserveDoesNotOverspend(t *testing.T) {
	ledger := capacity.NewLedger(100)

	var allowed atomic.Uint64
	var denied atomic.Uint64
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, _, ok := ledger.TryReserve(1)
			if ok {
				allowed.Add(1)
			} else {
				denied.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got != 100 {
		t.Fatalf("allowed reservations = %d, want 100", got)
	}
	if got := denied.Load(); got != 900 {
		t.Fatalf("denied reservations = %d, want 900", got)
	}
	requireSnapshotValue(t, ledger.Snapshot(), 100, 100, 0, 0)
}

func TestLedgerConcurrentReleaseRestoresCapacity(t *testing.T) {
	ledger := capacity.NewLedger(100)
	reservations := make([]*capacity.Reservation, 0, 100)
	for i := 0; i < 100; i++ {
		reservation, _, ok := ledger.TryReserve(1)
		if !ok {
			t.Fatalf("reservation %d failed", i)
		}
		reservations = append(reservations, reservation)
	}

	var wg sync.WaitGroup
	for _, reservation := range reservations {
		wg.Add(1)
		go func(res *capacity.Reservation) {
			defer wg.Done()
			res.Release()
		}(reservation)
	}
	wg.Wait()

	requireSnapshotValue(t, ledger.Snapshot(), 100, 0, 100, 0)
}

func TestLedgerConcurrentReserveReleaseAndSnapshotIsRaceFree(t *testing.T) {
	ledger := capacity.NewLedger(32)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				reservation, _, ok := ledger.TryReserve(1)
				if ok {
					_, _ = reservation.TryRelease()
				}
				_ = ledger.Snapshot()
				_ = ledger.Revision()
			}
		}()
	}
	wg.Wait()

	if snap := ledger.Snapshot(); !snap.Value.IsValid() {
		t.Fatalf("final snapshot is invalid: %+v", snap.Value)
	}
}

func TestReservationConcurrentTryReleaseReleasesExactlyOnce(t *testing.T) {
	ledger := capacity.NewLedger(1)
	reservation, _, ok := ledger.TryReserve(1)
	if !ok {
		t.Fatal("reservation failed")
	}

	var released atomic.Uint64
	var skipped atomic.Uint64
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			_, ok := reservation.TryRelease()
			if ok {
				released.Add(1)
			} else {
				skipped.Add(1)
			}
		}()
	}
	wg.Wait()

	if got := released.Load(); got != 1 {
		t.Fatalf("successful releases = %d, want 1", got)
	}
	if got := skipped.Load(); got != 99 {
		t.Fatalf("skipped releases = %d, want 99", got)
	}
	requireSnapshotValue(t, ledger.Snapshot(), 1, 0, 1, 0)
}

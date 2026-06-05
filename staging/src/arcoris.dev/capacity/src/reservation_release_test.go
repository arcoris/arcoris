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

func TestReservationReleaseReturnsDemand(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 4)))
	reservation := ledger.TryReserve(demand(t, entry("worker_slots", 3))).Reservation

	released := reservation.Release()
	if !reservation.Released() {
		t.Fatal("reservation not marked released")
	}
	requireVector(t, released.Value.Reserved)
	requireEntries(t, reservation.Demand().Entries(), entry("worker_slots", 3))
}

func TestReservationTryReleaseIdempotent(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 4)))
	reservation := ledger.TryReserve(demand(t, entry("worker_slots", 3))).Reservation
	first, ok := reservation.TryRelease()
	if !ok {
		t.Fatal("first TryRelease() returned false")
	}
	second, ok := reservation.TryRelease()
	if ok {
		t.Fatal("second TryRelease() returned true")
	}
	if second.Revision != first.Revision {
		t.Fatal("second TryRelease advanced revision")
	}
}

func TestReservationReleasePanicsOnDoubleRelease(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 4)))
	reservation := ledger.TryReserve(demand(t, entry("worker_slots", 3))).Reservation
	_ = reservation.Release()
	requirePanicIs(t, capacity.ErrReservationReleased, func() { _ = reservation.Release() })
}

func TestReservationConcurrentTryReleaseReleasesOnce(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 1)))
	reservation := ledger.TryReserve(demand(t, entry("worker_slots", 1))).Reservation

	var wg sync.WaitGroup
	results := make(chan bool, 8)
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, ok := reservation.TryRelease()
			results <- ok
		}()
	}
	wg.Wait()
	close(results)

	successes := 0
	for ok := range results {
		if ok {
			successes++
		}
	}
	if successes != 1 {
		t.Fatalf("successful releases = %d, want 1", successes)
	}
}

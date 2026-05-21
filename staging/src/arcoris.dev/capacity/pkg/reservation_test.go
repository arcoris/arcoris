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
	"testing"

	"arcoris.dev/capacity"
)

func TestReservationReleaseReturnsCapacity(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}
	if reservation.Released() {
		t.Fatal("reservation is released before Release")
	}

	snap := reservation.Release()
	requireSnapshotValue(t, snap, 10, 0, 10, 0)
	if !reservation.Released() {
		t.Fatal("reservation is not released after Release")
	}
}

func TestReservationAmountRemainsStableAfterRelease(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}
	if got := reservation.Amount(); got != 4 {
		t.Fatalf("Amount() before release = %d, want 4", got)
	}

	reservation.Release()
	if got := reservation.Amount(); got != 4 {
		t.Fatalf("Amount() after release = %d, want 4", got)
	}
}

func TestReservationReleaseAdvancesRevision(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, before, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}
	after := reservation.Release()

	if !after.Revision.ChangedSince(before.Revision) {
		t.Fatalf("revision did not advance after Release")
	}
}

func TestReservationReleaseTwicePanics(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}
	reservation.Release()

	requirePanic(t, "capacity.Reservation: already released", func() { reservation.Release() })
}

func TestReservationTryReleaseIsIdempotent(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}

	first, ok := reservation.TryRelease()
	if !ok {
		t.Fatal("first TryRelease returned ok=false, want true")
	}
	requireSnapshotValue(t, first, 10, 0, 10, 0)

	second, ok := reservation.TryRelease()
	if ok {
		t.Fatal("second TryRelease returned ok=true, want false")
	}
	if second.Revision != first.Revision {
		t.Fatalf("second revision = %d, want unchanged %d", second.Revision, first.Revision)
	}
	requireSnapshotValue(t, second, 10, 0, 10, 0)
}

func TestNilReservationPanics(t *testing.T) {
	t.Parallel()

	var reservation *capacity.Reservation
	requirePanic(t, "capacity.Reservation: nil reservation", func() { _ = reservation.Amount() })
	requirePanic(t, "capacity.Reservation: nil reservation", func() { _ = reservation.Released() })
	requirePanic(t, "capacity.Reservation: nil reservation", func() { _ = reservation.Release() })
	requirePanic(t, "capacity.Reservation: nil reservation", func() { _, _ = reservation.TryRelease() })
}

func TestZeroReservationPanics(t *testing.T) {
	t.Parallel()

	var reservation capacity.Reservation
	requirePanic(t, "capacity.Reservation: invalid reservation", func() { _ = reservation.Amount() })
	requirePanic(t, "capacity.Reservation: invalid reservation", func() { _ = reservation.Released() })
	requirePanic(t, "capacity.Reservation: invalid reservation", func() { _ = reservation.Release() })
	requirePanic(t, "capacity.Reservation: invalid reservation", func() { _, _ = reservation.TryRelease() })
}

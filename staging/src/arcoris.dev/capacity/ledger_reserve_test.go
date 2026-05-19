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

func TestLedgerTryReserveSucceedsWhenCapacityIsAvailable(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	initialRevision := ledger.Revision()

	reservation, snap, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("TryReserve returned ok=false, want true")
	}
	if reservation == nil {
		t.Fatal("reservation is nil")
	}
	if reservation.Amount() != 4 {
		t.Fatalf("reservation amount = %d, want 4", reservation.Amount())
	}
	if !snap.Revision.ChangedSince(initialRevision) {
		t.Fatalf("revision did not advance after successful reservation")
	}
	requireSnapshotValue(t, snap, 10, 4, 6, 0)
}

func TestLedgerTryReserveSucceedsAtExactLimit(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, snap, ok := ledger.TryReserve(10)
	if !ok {
		t.Fatal("TryReserve exact limit returned ok=false, want true")
	}
	if reservation == nil {
		t.Fatal("reservation is nil")
	}
	requireSnapshotValue(t, snap, 10, 10, 0, 0)
}

func TestLedgerTryReserveFailsWhenCapacityIsInsufficient(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	_, snap, ok := ledger.TryReserve(7)
	if !ok {
		t.Fatal("initial reservation failed")
	}
	beforeDenied := snap.Revision

	reservation, denied, ok := ledger.TryReserve(4)
	if ok {
		t.Fatal("TryReserve returned ok=true, want false")
	}
	if reservation != nil {
		t.Fatalf("reservation = %#v, want nil", reservation)
	}
	if denied.Revision != beforeDenied {
		t.Fatalf("denied revision = %d, want unchanged %d", denied.Revision, beforeDenied)
	}
	requireSnapshotValue(t, denied, 10, 7, 3, 0)
	current := ledger.Snapshot()
	if denied != current {
		t.Fatalf("denied snapshot = %+v, want current snapshot %+v", denied, current)
	}
}

func TestLedgerTryReserveFailsWhenLimitIsZero(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(0)
	reservation, snap, ok := ledger.TryReserve(1)
	if ok {
		t.Fatal("TryReserve returned ok=true, want false")
	}
	if reservation != nil {
		t.Fatalf("reservation = %#v, want nil", reservation)
	}
	requireSnapshotValue(t, snap, 0, 0, 0, 0)
}

func TestLedgerTryReserveZeroAmountPanics(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	requirePanic(t, "capacity: reservation amount must be positive", func() { _, _, _ = ledger.TryReserve(0) })
}

func TestLedgerTryReserveFailsWhileOvercommitted(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	_, _, ok := ledger.TryReserve(8)
	if !ok {
		t.Fatal("initial reservation failed")
	}
	ledger.SetLimit(5)

	reservation, snap, ok := ledger.TryReserve(1)
	if ok {
		t.Fatal("TryReserve returned ok=true while overcommitted")
	}
	if reservation != nil {
		t.Fatalf("reservation = %#v, want nil", reservation)
	}
	requireSnapshotValue(t, snap, 5, 8, 0, 3)
}

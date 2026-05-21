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

func TestLedgerSetLimitIncreasesAvailableCapacity(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	_, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}

	snap := ledger.SetLimit(15)
	requireSnapshotValue(t, snap, 15, 4, 11, 0)
}

func TestLedgerSetLimitDecreasesAvailableCapacity(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	_, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}

	snap := ledger.SetLimit(7)
	requireSnapshotValue(t, snap, 7, 4, 3, 0)
}

func TestLedgerSetLimitBelowReservedCreatesDebt(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, _, ok := ledger.TryReserve(8)
	if !ok {
		t.Fatal("reservation failed")
	}

	snap := ledger.SetLimit(5)
	requireSnapshotValue(t, snap, 5, 8, 0, 3)
	if reservation.Released() {
		t.Fatal("reservation was revoked by limit reduction")
	}
}

func TestLedgerSetLimitSameValueDoesNotAdvanceRevision(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	before := ledger.Revision()
	snap := ledger.SetLimit(10)

	if snap.Revision != before {
		t.Fatalf("revision = %d, want unchanged %d", snap.Revision, before)
	}
	requireSnapshotValue(t, snap, 10, 0, 10, 0)
}

func TestLedgerSetLimitToZeroKeepsExistingReservations(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	reservation, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}

	snap := ledger.SetLimit(0)
	requireSnapshotValue(t, snap, 0, 4, 0, 4)
	if reservation.Released() {
		t.Fatal("reservation was revoked by SetLimit(0)")
	}
}

func TestLedgerReleaseReducesDebtAndRestoresAvailability(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	large, _, ok := ledger.TryReserve(6)
	if !ok {
		t.Fatal("large reservation failed")
	}
	small, _, ok := ledger.TryReserve(2)
	if !ok {
		t.Fatal("small reservation failed")
	}

	ledger.SetLimit(5)
	debtAfterSmallRelease := small.Release()
	requireSnapshotValue(t, debtAfterSmallRelease, 5, 6, 0, 1)

	availableAfterLargeRelease := large.Release()
	requireSnapshotValue(t, availableAfterLargeRelease, 5, 0, 5, 0)
}

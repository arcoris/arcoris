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
	"arcoris.dev/snapshot"
	panicassert "arcoris.dev/testutil/panic"
)

func TestLedgerImplementsSnapshotSources(t *testing.T) {
	t.Parallel()

	var _ snapshot.Source[capacity.Snapshot] = (*capacity.Ledger)(nil)
	var _ snapshot.RevisionSource = (*capacity.Ledger)(nil)
}

func TestNewLedgerInitialSnapshot(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(10)
	snap := ledger.Snapshot()
	rev := ledger.Revision()

	if snap.Revision != rev {
		t.Fatalf("Snapshot revision = %d, Revision() = %d", snap.Revision, rev)
	}
	requireSnapshotValue(t, snap, 10, 0, 10, 0)
}

func TestNewLedgerAllowsZeroLimit(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(0)
	snap := ledger.Snapshot()

	requireSnapshotValue(t, snap, 0, 0, 0, 0)
	if !snap.Value.Exhausted() {
		t.Fatalf("zero-limit snapshot should be exhausted: %+v", snap.Value)
	}
}

func TestZeroLedgerPanics(t *testing.T) {
	t.Parallel()

	var ledger capacity.Ledger
	panicassert.RequireMessage(t, "capacity.Ledger: uninitialized ledger", func() { _ = ledger.Snapshot() })
	panicassert.RequireMessage(t, "capacity.Ledger: uninitialized ledger", func() { _ = ledger.Revision() })
	panicassert.RequireMessage(t, "capacity.Ledger: uninitialized ledger", func() { _ = ledger.SetLimit(1) })
	panicassert.RequireMessage(t, "capacity.Ledger: uninitialized ledger", func() { _, _, _ = ledger.TryReserve(1) })
	panicassert.RequireMessage(t, "capacity.Ledger: uninitialized ledger", func() { _, _, _ = ledger.TryReserve(0) })
}

func TestNilLedgerPanics(t *testing.T) {
	t.Parallel()

	var ledger *capacity.Ledger
	panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() { _ = ledger.Snapshot() })
	panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() { _ = ledger.Revision() })
	panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() { _ = ledger.SetLimit(1) })
	panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() { _, _, _ = ledger.TryReserve(1) })
	panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() { _, _, _ = ledger.TryReserve(0) })
}

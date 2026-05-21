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


package bulkhead

import "testing"

func TestSetLimitBelowActiveLeasesCreatesDebt(t *testing.T) {
	t.Parallel()

	b := New(3)
	first, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("first TryAcquire failed")
	}
	defer first.TryRelease()
	second, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("second TryAcquire failed")
	}
	defer second.TryRelease()

	before := b.Revision()
	snap := b.SetLimit(1)
	if snap.Revision == before {
		t.Fatal("SetLimit did not advance revision")
	}
	requireSnapshotValue(t, snap, 1, 2, 0, 1)
}

func TestSetLimitSameValueDoesNotAdvanceRevision(t *testing.T) {
	t.Parallel()

	b := New(2)
	before := b.Revision()
	snap := b.SetLimit(2)
	if snap.Revision != before {
		t.Fatalf("revision = %d, want unchanged %d", snap.Revision, before)
	}
	requireSnapshotValue(t, snap, 2, 0, 2, 0)
}

func TestReleaseAfterDebtRestoresAvailability(t *testing.T) {
	t.Parallel()

	b := New(3)
	first, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("first TryAcquire failed")
	}
	second, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("second TryAcquire failed")
	}

	b.SetLimit(1)
	debt := first.Release()
	requireSnapshotValue(t, debt, 1, 1, 0, 0)

	available := second.Release()
	requireSnapshotValue(t, available, 1, 0, 1, 0)
}

func TestSetLimitRaisesClosedBulkhead(t *testing.T) {
	t.Parallel()

	b := New(0)
	opened := b.SetLimit(1)
	requireSnapshotValue(t, opened, 1, 0, 1, 0)

	lease, snap, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire after SetLimit failed")
	}
	defer lease.Release()
	requireSnapshotValue(t, snap, 1, 1, 0, 0)
}

func TestSetLimitToZeroKeepsActiveLeases(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}
	defer lease.TryRelease()

	closed := b.SetLimit(0)
	requireSnapshotValue(t, closed, 0, 1, 0, 1)
	if lease.Released() {
		t.Fatal("active lease was revoked by SetLimit(0)")
	}
}

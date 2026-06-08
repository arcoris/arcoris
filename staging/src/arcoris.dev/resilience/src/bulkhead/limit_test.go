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

func TestTryAcquireDeniedWhileDebtExists(t *testing.T) {
	t.Parallel()

	b := New(2)
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

	b.SetLimit(1)
	lease, observation, ok := b.TryAcquire()
	if ok {
		t.Fatal("TryAcquire under debt returned ok=true, want false")
	}
	if lease != nil {
		t.Fatalf("denied lease = %#v, want nil", lease)
	}
	requireObservationValue(t, observation, RefusalDebt, 1, 2, 0, 1)
}

func TestTryAcquireAmountDeniedWhileDebtExists(t *testing.T) {
	t.Parallel()

	b := New(2)
	held, _, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("TryAcquireAmount failed")
	}
	defer held.TryRelease()

	b.SetLimit(1)
	lease, observation, ok := b.TryAcquireAmount(1)
	if ok {
		t.Fatal("TryAcquireAmount under debt returned ok=true, want false")
	}
	if lease != nil {
		t.Fatalf("denied lease = %#v, want nil", lease)
	}
	requireObservationValue(t, observation, RefusalDebt, 1, 2, 0, 1)
}

func TestIncreasingLimitWhileDebtExistsRestoresAvailability(t *testing.T) {
	t.Parallel()

	b := New(2)
	held, _, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("TryAcquireAmount failed")
	}
	defer held.TryRelease()

	b.SetLimit(1)
	opened := b.SetLimit(3)
	requireSnapshotValue(t, opened, 3, 2, 1, 0)

	next, observation, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire after limit increase failed")
	}
	defer next.TryRelease()
	requireObservationValue(t, observation, RefusalNone, 3, 3, 0, 0)
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

	lease, observation, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire after SetLimit failed")
	}
	defer lease.Release()
	requireObservationValue(t, observation, RefusalNone, 1, 1, 0, 0)
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

func TestBulkheadLeaseReleaseAfterLimitIncrease(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}

	b.SetLimit(2)
	if lease.Released() {
		t.Fatal("lease was released by limit increase")
	}

	snap := lease.Release()
	requireSnapshotValue(t, snap, 2, 0, 2, 0)
}

func TestBulkheadLeaseReleaseAfterLimitDecrease(t *testing.T) {
	t.Parallel()

	b := New(3)
	lease, _, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("TryAcquireAmount failed")
	}

	b.SetLimit(1)
	if lease.Released() {
		t.Fatal("lease was released by limit decrease")
	}

	snap := lease.Release()
	requireSnapshotValue(t, snap, 1, 0, 1, 0)
}

func TestBulkheadLeaseReleaseAfterLimitSetToZero(t *testing.T) {
	t.Parallel()

	b := New(2)
	first, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("first TryAcquire failed")
	}
	second, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("second TryAcquire failed")
	}

	b.SetLimit(0)
	if first.Released() || second.Released() {
		t.Fatal("lease was released by SetLimit(0)")
	}

	snap := first.Release()
	requireSnapshotValue(t, snap, 0, 1, 0, 1)
	snap = second.Release()
	requireSnapshotValue(t, snap, 0, 0, 0, 0)
}

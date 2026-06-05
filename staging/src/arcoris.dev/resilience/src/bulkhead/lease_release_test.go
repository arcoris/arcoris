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

import (
	"testing"

	"arcoris.dev/capacity"
	panicassert "arcoris.dev/testutil/panic"
)

func TestLeaseReleaseRestoresAvailableCapacity(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}

	before := b.Revision()
	snap := lease.Release()
	if snap.Revision == before {
		t.Fatal("Release did not advance revision")
	}
	requireSnapshotValue(t, snap, 1, 0, 1, 0)
}

func TestLeaseTryReleaseIsIdempotent(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}

	first, ok := lease.TryRelease()
	if !ok {
		t.Fatal("first TryRelease returned ok=false, want true")
	}
	requireSnapshotValue(t, first, 1, 0, 1, 0)

	second, ok := lease.TryRelease()
	if ok {
		t.Fatal("second TryRelease returned ok=true, want false")
	}
	if second.Revision != first.Revision {
		t.Fatalf("second revision = %d, want unchanged %d", second.Revision, first.Revision)
	}
	requireSnapshotValue(t, second, 1, 0, 1, 0)
}

func TestLeaseReleasePanicsAfterTryRelease(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}
	_, _ = lease.TryRelease()

	panicassert.RequireErrorIs(t, capacity.ErrReservationReleased, func() {
		_ = lease.Release()
	})
}

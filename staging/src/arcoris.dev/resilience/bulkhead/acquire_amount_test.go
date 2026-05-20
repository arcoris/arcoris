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

package bulkhead

import "testing"

func TestTryAcquireAmountReservesWeightedCapacity(t *testing.T) {
	t.Parallel()

	b := New(3)
	lease, snap, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("TryAcquireAmount returned ok=false, want true")
	}
	defer lease.Release()
	if lease == nil {
		t.Fatal("lease is nil")
	}
	if lease.Amount() != 2 {
		t.Fatalf("lease amount = %d, want 2", lease.Amount())
	}
	requireSnapshotValue(t, snap, 3, 2, 1, 0)
}

func TestTryAcquireAmountOneMatchesTryAcquireBehavior(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, snap, ok := b.TryAcquireAmount(1)
	if !ok {
		t.Fatal("TryAcquireAmount(1) returned ok=false, want true")
	}
	defer lease.Release()
	requireSnapshotValue(t, snap, 1, 1, 0, 0)

	deniedLease, deniedSnap, ok := b.TryAcquire()
	if ok {
		t.Fatal("TryAcquire returned ok=true, want false")
	}
	if deniedLease != nil {
		t.Fatalf("denied lease = %#v, want nil", deniedLease)
	}
	requireSnapshotValue(t, deniedSnap, 1, 1, 0, 0)
}

func TestTryAcquireAmountDeniedWhenCapacityIsInsufficient(t *testing.T) {
	t.Parallel()

	b := New(3)
	lease, _, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("first TryAcquireAmount failed")
	}
	defer lease.Release()

	deniedLease, snap, ok := b.TryAcquireAmount(2)
	if ok {
		t.Fatal("second TryAcquireAmount returned ok=true, want false")
	}
	if deniedLease != nil {
		t.Fatalf("denied lease = %#v, want nil", deniedLease)
	}
	requireSnapshotValue(t, snap, 3, 2, 1, 0)
}

func TestTryAcquireAmountReleaseRestoresCapacity(t *testing.T) {
	t.Parallel()

	b := New(3)
	lease, _, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("TryAcquireAmount failed")
	}
	lease.Release()

	next, snap, ok := b.TryAcquireAmount(3)
	if !ok {
		t.Fatal("TryAcquireAmount after release returned ok=false, want true")
	}
	defer next.Release()
	requireSnapshotValue(t, snap, 3, 3, 0, 0)
}

func TestTryAcquireAmountInvalidAmountPanics(t *testing.T) {
	t.Parallel()

	b := New(1)
	requirePanic(t, "capacity: reservation amount must be positive", func() {
		_, _, _ = b.TryAcquireAmount(0)
	})
}

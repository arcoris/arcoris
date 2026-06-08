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

func TestTryAcquireSucceedsUnderLimit(t *testing.T) {
	t.Parallel()

	b := New(2)
	before := b.Revision()
	lease, observation, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire returned ok=false, want true")
	}
	defer lease.Release()
	if lease == nil {
		t.Fatal("lease is nil")
	}
	if lease.Amount() != 1 {
		t.Fatalf("lease amount = %d, want 1", lease.Amount())
	}
	if observation.Refusal != RefusalNone {
		t.Fatalf("refusal = %s, want none", observation.Refusal)
	}
	if observation.Snapshot.Revision == before {
		t.Fatal("successful acquire did not advance revision")
	}
	requireObservationValue(t, observation, RefusalNone, 2, 1, 1, 0)
}

func TestTryAcquireDeniedWhenCapacityExhausted(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("first TryAcquire failed")
	}
	defer lease.Release()

	beforeDenied := b.Revision()
	deniedLease, observation, ok := b.TryAcquire()
	if ok {
		t.Fatal("second TryAcquire returned ok=true, want false")
	}
	if deniedLease != nil {
		t.Fatalf("denied lease = %#v, want nil", deniedLease)
	}
	requireObservationValue(t, observation, RefusalInsufficient, 1, 1, 0, 0)
	if got := b.Snapshot(); got != observation.Snapshot {
		t.Fatalf("denied snapshot = %+v, want current snapshot %+v", observation.Snapshot, got)
	}
	if got := b.Revision(); got != beforeDenied {
		t.Fatalf("denied acquire advanced revision: got %d, want %d", got, beforeDenied)
	}
}

func TestTryAcquireDeniedWhenLimitIsZero(t *testing.T) {
	t.Parallel()

	b := New(0)
	lease, observation, ok := b.TryAcquire()
	if ok {
		t.Fatal("TryAcquire returned ok=true, want false")
	}
	if lease != nil {
		t.Fatalf("lease = %#v, want nil", lease)
	}
	requireObservationValue(t, observation, RefusalInsufficient, 0, 0, 0, 0)
}

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

import (
	"errors"
	"testing"
)

func TestNewPublishesInitialSnapshot(t *testing.T) {
	t.Parallel()

	l, err := New(3)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	snap := l.Snapshot()
	if snap.IsZeroRevision() {
		t.Fatal("initial revision is zero")
	}
	if !snap.Value.IsValid() {
		t.Fatalf("initial snapshot is invalid: %+v", snap.Value)
	}
	if snap.Value.Capacity.Limit != 3 {
		t.Fatalf("Limit = %d, want 3", snap.Value.Capacity.Limit)
	}
	if snap.Value.Capacity.InFlight != 0 || snap.Value.Capacity.Available != 3 || snap.Value.Capacity.Full {
		t.Fatalf("unexpected initial capacity: %+v", snap.Value.Capacity)
	}
	if snap.Value.Stats != (StatsSnapshot{}) {
		t.Fatalf("initial stats = %+v, want zero", snap.Value.Stats)
	}
}

func TestNewRejectsInvalidLimit(t *testing.T) {
	t.Parallel()

	l, err := New(0)
	if l != nil {
		t.Fatalf("limiter = %v, want nil", l)
	}
	if !errors.Is(err, ErrInvalidLimit) {
		t.Fatalf("New error = %v, want %v", err, ErrInvalidLimit)
	}
}

func TestLimiterTryAcquireAllowed(t *testing.T) {
	t.Parallel()

	l, err := New(2)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	prev := l.Revision()
	permit, dec := l.TryAcquire()
	if permit == nil {
		t.Fatal("permit is nil")
	}
	defer permit.Release()
	if !dec.Allowed || dec.Reason != ReasonAllowed {
		t.Fatalf("decision = %+v, want allowed", dec)
	}
	if !dec.IsValid() {
		t.Fatalf("decision is invalid: %+v", dec)
	}
	if dec.Snapshot.Revision == prev {
		t.Fatalf("revision did not advance: %s", dec.Snapshot.Revision)
	}
	if got := dec.Snapshot.Value.Capacity.InFlight; got != 1 {
		t.Fatalf("InFlight = %d, want 1", got)
	}
	if got := dec.Snapshot.Value.Stats.Acquired; got != 1 {
		t.Fatalf("Acquired = %d, want 1", got)
	}
}

func TestLimiterTryAcquireDeniedWhenFull(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	permit, dec := l.TryAcquire()
	if permit == nil || !dec.Allowed {
		t.Fatalf("first acquire = permit %v decision %+v, want allowed", permit, dec)
	}
	defer permit.Release()

	prev := l.Revision()
	permit, denied := l.TryAcquire()
	if permit != nil {
		t.Fatalf("denied permit = %v, want nil", permit)
	}
	if denied.Allowed || denied.Reason != ReasonFull {
		t.Fatalf("decision = %+v, want full", denied)
	}
	if !denied.IsValid() {
		t.Fatalf("denied decision is invalid: %+v", denied)
	}
	if denied.Snapshot.Revision == prev {
		t.Fatal("rejected acquire did not advance revision")
	}
	if got := denied.Snapshot.Value.Stats.Rejected; got != 1 {
		t.Fatalf("Rejected = %d, want 1", got)
	}
	if got := denied.Snapshot.Value.Capacity.InFlight; got != 1 {
		t.Fatalf("InFlight = %d, want 1", got)
	}
}

func TestLimiterSnapshotStampedAndRevision(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	snap := l.Snapshot()
	stamped := l.Stamped()
	if stamped.Revision != snap.Revision {
		t.Fatalf("Stamped revision = %s, want %s", stamped.Revision, snap.Revision)
	}
	if l.Revision() != snap.Revision {
		t.Fatalf("Revision() = %s, want %s", l.Revision(), snap.Revision)
	}
	if !stamped.Value.IsValid() {
		t.Fatalf("stamped value invalid: %+v", stamped.Value)
	}
}

func TestLimiterPanicsOnNilOrUninitializedReceiver(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		defer func() {
			if got := recover(); !errors.Is(asError(got), ErrNilLimiter) {
				t.Fatalf("panic = %v, want %v", got, ErrNilLimiter)
			}
		}()

		var l *Limiter
		_ = l.Snapshot()
	})

	t.Run("uninitialized", func(t *testing.T) {
		defer func() {
			if got := recover(); !errors.Is(asError(got), ErrUninitializedLimiter) {
				t.Fatalf("panic = %v, want %v", got, ErrUninitializedLimiter)
			}
		}()

		var l Limiter
		_ = l.Snapshot()
	})
}

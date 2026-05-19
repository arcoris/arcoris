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

func TestPermitReleaseReturnsCapacity(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	permit, dec := l.TryAcquire()
	if permit == nil || !dec.Allowed {
		t.Fatalf("TryAcquire = %v %+v, want allowed", permit, dec)
	}

	prev := l.Revision()
	permit.Release()
	if !permit.Released() {
		t.Fatal("Released() = false, want true")
	}

	snap := l.Snapshot()
	if snap.Revision == prev {
		t.Fatal("release did not advance revision")
	}
	if got := snap.Value.Capacity.InFlight; got != 0 {
		t.Fatalf("InFlight = %d, want 0", got)
	}
	if got := snap.Value.Stats.Released; got != 1 {
		t.Fatalf("Released = %d, want 1", got)
	}

	permit2, dec2 := l.TryAcquire()
	if permit2 == nil || !dec2.Allowed {
		t.Fatalf("second TryAcquire = %v %+v, want allowed", permit2, dec2)
	}
	permit2.Release()
}

func TestPermitReleaseIsIdempotent(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	permit, _ := l.TryAcquire()
	permit.Release()
	first := l.Snapshot()
	permit.Release()
	second := l.Snapshot()

	if first.Revision != second.Revision {
		t.Fatalf("double release advanced revision: first %s second %s", first.Revision, second.Revision)
	}
	if second.Value.Stats.Released != 1 {
		t.Fatalf("Released = %d, want 1", second.Value.Stats.Released)
	}
}

func TestNilPermitReleaseIsNoop(t *testing.T) {
	t.Parallel()

	var p *Permit
	p.Release()
	if !p.Released() {
		t.Fatal("nil permit Released() = false, want true")
	}
}

func TestReleaseUnderflowPanics(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	defer func() {
		if got := recover(); !errors.Is(asError(got), ErrReleaseUnderflow) {
			t.Fatalf("panic = %v, want %v", got, ErrReleaseUnderflow)
		}
	}()

	l.release()
}

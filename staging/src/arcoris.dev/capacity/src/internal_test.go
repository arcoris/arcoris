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


package capacity

import "testing"

func TestReservationTryReleaseDetectsCorruptedReservedAccounting(t *testing.T) {
	ledger := NewLedger(10)
	reservation, _, ok := ledger.TryReserve(4)
	if !ok {
		t.Fatal("reservation failed")
	}

	// Intentionally corrupt private ledger state to verify the defensive
	// underflow invariant panic. Public API calls do not expose this state.
	ledger.mu.Lock()
	ledger.reserved = 3
	ledger.mu.Unlock()

	requireInternalPanic(t, errLedgerReservedUnderflow, func() {
		_, _ = reservation.TryRelease()
	})
}

func TestReservationRequireNonNilPanicsForDetachedReservation(t *testing.T) {
	reservation := &Reservation{}

	requireInternalPanic(t, errInvalidReservation, func() {
		reservation.requireNonNil()
	})
}

func requireInternalPanic(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %q", want)
		}
		if got, ok := recovered.(string); !ok || got != want {
			t.Fatalf("panic = %#v, want %q", recovered, want)
		}
	}()

	fn()
}

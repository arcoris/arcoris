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
)

func TestVectorReservationReleaseObserved(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t, entry("worker_slots", 4)))
	reservation, ok := ledger.TryReserve(demand(t, entry("worker_slots", 3)))
	if !ok {
		t.Fatal("TryReserve() failed")
	}

	released := reservation.ReleaseObserved()
	if !released.Value.Reserved.IsZero() {
		t.Fatalf("reserved after release = %#v", released.Value.Reserved)
	}
	requirePanicIs(t, capacity.ErrReservationReleased, reservation.Release)
}

func TestVectorReservationTryReleaseObserved(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t, entry("worker_slots", 2)))
	reservation, ok := ledger.TryReserve(demand(t, entry("worker_slots", 1)))
	if !ok {
		t.Fatal("TryReserve() failed")
	}

	first, ok := reservation.TryReleaseObserved()
	if !ok || !first.Value.Reserved.IsZero() {
		t.Fatalf("first TryReleaseObserved() = %#v, %v", first.Value, ok)
	}

	second, ok := reservation.TryReleaseObserved()
	if ok || !second.Value.Reserved.IsZero() {
		t.Fatalf("second TryReleaseObserved() = %#v, %v", second.Value, ok)
	}
}

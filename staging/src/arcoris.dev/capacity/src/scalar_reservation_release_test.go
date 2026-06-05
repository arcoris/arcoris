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

func TestScalarReservationReleaseRestoresCapacity(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(4)
	result := ledger.TryReserve(3)

	released := result.Reservation.Release()
	if released.Value.Reserved != 0 || released.Value.Available != 4 {
		t.Fatalf("release snapshot = %+v, want reserved=0 available=4", released.Value)
	}
}

func TestScalarReservationTryReleaseIdempotent(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(2)
	reservation := ledger.TryReserve(1).Reservation
	first, ok := reservation.TryRelease()
	if !ok {
		t.Fatal("first TryRelease() returned false")
	}
	second, ok := reservation.TryRelease()
	if ok {
		t.Fatal("second TryRelease() returned true")
	}
	if second.Revision != first.Revision {
		t.Fatal("second TryRelease advanced revision")
	}
	requirePanicIs(t, capacity.ErrReservationReleased, func() { _ = reservation.Release() })
}

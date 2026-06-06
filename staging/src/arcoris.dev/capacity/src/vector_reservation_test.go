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

func TestVectorReservationDemandAndReleased(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t, entry("worker_slots", 4)))
	reservation, ok := ledger.TryReserve(demand(t, entry("worker_slots", 3)))
	if !ok {
		t.Fatal("TryReserve() failed")
	}

	requireVector(t, reservation.Demand().Vector(), entry("worker_slots", 3))
	if reservation.Released() {
		t.Fatal("Released() before release = true")
	}

	reservation.Release()

	if !reservation.Released() {
		t.Fatal("Released() after release = false")
	}
}

func TestVectorReservationValidatePanics(t *testing.T) {
	var nilReservation *capacity.VectorReservation
	requirePanicIs(t, capacity.ErrNilReservation, func() { _ = nilReservation.Demand() })

	var zeroReservation capacity.VectorReservation
	requirePanicIs(t, capacity.ErrInvalidReservation, func() { _ = zeroReservation.Demand() })
}

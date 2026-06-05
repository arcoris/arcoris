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

func TestReservationDemandAndReleased(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewLedger(vector(t, entry("worker_slots", 4)))
	result := ledger.TryReserve(demand(t, entry("worker_slots", 3)))
	reservation := result.Reservation
	if reservation == nil {
		t.Fatal("reservation is nil")
	}
	requireEntries(t, reservation.Demand().Entries(), entry("worker_slots", 3))
	if reservation.Released() {
		t.Fatal("new reservation is released")
	}
}

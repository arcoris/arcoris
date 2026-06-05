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

// ScalarReservation owns one amount reserved from a ScalarLedger.
type ScalarReservation struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// ledger is the owner that created this reservation.
	ledger *ScalarLedger

	// amount is the immutable scalar amount owned by this reservation.
	amount Amount

	// released is protected by ledger.mu.
	released bool
}

// Amount returns the amount reserved by r.
func (r *ScalarReservation) Amount() Amount {
	r.requireNonNil()
	return r.amount
}

// Released reports whether r has already been released.
func (r *ScalarReservation) Released() bool {
	r.requireNonNil()

	r.ledger.mu.Lock()
	defer r.ledger.mu.Unlock()

	r.ledger.requireInitializedLocked()
	return r.released
}

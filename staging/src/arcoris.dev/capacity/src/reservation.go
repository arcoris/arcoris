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

import "sync/atomic"

// Reservation owns scalar capacity acquired from a Ledger.
//
// Reservation is returned only by a successful Ledger.TryAcquire or
// Ledger.TryAcquireObserved call. It must not be copied after creation. Amount
// remains observable after release for cleanup diagnostics.
type Reservation struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// ledger is the owner that created this reservation.
	ledger *Ledger

	// amount is the immutable scalar amount owned by this reservation.
	amount Amount

	// released records whether this reservation has already returned amount.
	released atomic.Bool
}

// Amount returns the amount reserved by r.
func (r *Reservation) Amount() Amount {
	r.requireValid()

	return r.amount
}

// Released reports whether r has already been released.
func (r *Reservation) Released() bool {
	r.requireValid()

	return r.released.Load()
}

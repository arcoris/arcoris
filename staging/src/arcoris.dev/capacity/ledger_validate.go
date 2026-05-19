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

package capacity

// requirePositiveAmount panics when amount cannot represent a reservation
// request.
func requirePositiveAmount(amount Amount) {
	if amount == 0 {
		panic(errZeroReservationAmount)
	}
}

// requireNonNil panics when l is nil.
func (l *Ledger) requireNonNil() {
	if l == nil {
		panic(errNilLedger)
	}
}

// requireInitializedLocked panics when l is a zero Ledger value.
//
// The caller must hold l.mu. Initialization is checked under the ledger lock so
// ordinary concurrent operations observe a stable revision field while they
// validate ledger ownership.
func (l *Ledger) requireInitializedLocked() {
	if l.revision.IsZero() {
		panic(errUninitializedLedger)
	}
}

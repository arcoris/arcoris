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

const (
	// errNilLedger is the panic value used when a method is called on a nil
	// Ledger.
	errNilLedger = "capacity.Ledger: nil ledger"

	// errUninitializedLedger is the panic value used when a method is called on a
	// zero Ledger value instead of a value created by NewLedger.
	errUninitializedLedger = "capacity.Ledger: uninitialized ledger"

	// errZeroReservationAmount is the panic value used when a caller attempts to
	// reserve zero capacity units.
	errZeroReservationAmount = "capacity: reservation amount must be positive"

	// errNilReservation is the panic value used when a method is called on a nil
	// Reservation.
	errNilReservation = "capacity.Reservation: nil reservation"

	// errInvalidReservation is the panic value used when a Reservation was not
	// created by a Ledger or has lost its owning ledger.
	errInvalidReservation = "capacity.Reservation: invalid reservation"

	// errReservationAlreadyReleased is the panic value used when Release is called
	// more than once for the same Reservation.
	errReservationAlreadyReleased = "capacity.Reservation: already released"

	// errLedgerReservedUnderflow is the panic value used when internal ledger
	// accounting would subtract more reserved capacity than exists.
	errLedgerReservedUnderflow = "capacity.Ledger: reserved capacity underflow"
)

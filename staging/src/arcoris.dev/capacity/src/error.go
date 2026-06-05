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

import (
	"errors"

	"arcoris.dev/capacity/internal/diagnostic"
)

// Sentinel errors identify broad capacity failure classes.
//
// Structured *Error values wrap these sentinels so callers can use errors.Is
// without depending on diagnostic strings.
var (
	// ErrInvalidResource identifies an invalid resource identity.
	ErrInvalidResource = errors.New("capacity: invalid resource")

	// ErrZeroAmount identifies a vector or demand entry with zero amount.
	ErrZeroAmount = errors.New("capacity: zero amount")

	// ErrDuplicateResource identifies repeated resources in constructor input.
	ErrDuplicateResource = errors.New("capacity: duplicate resource")

	// ErrEmptyDemand identifies an empty reservation demand.
	ErrEmptyDemand = errors.New("capacity: empty demand")

	// ErrAmountOverflow identifies checked amount addition overflow.
	ErrAmountOverflow = errors.New("capacity: amount overflow")

	// ErrAmountUnderflow identifies checked amount subtraction underflow.
	ErrAmountUnderflow = errors.New("capacity: amount underflow")

	// ErrInvalidVector identifies a vector that is not canonical.
	ErrInvalidVector = errors.New("capacity: invalid vector")

	// ErrInvalidDemand identifies a demand that is empty or not canonical.
	ErrInvalidDemand = errors.New("capacity: invalid demand")

	// ErrInvalidState identifies source accounting state that is not canonical.
	ErrInvalidState = errors.New("capacity: invalid state")

	// ErrNilLedger identifies a method call on a nil ledger receiver.
	ErrNilLedger = errors.New("capacity: nil ledger")

	// ErrUninitializedLedger identifies a zero-value stateful ledger owner.
	ErrUninitializedLedger = errors.New("capacity: uninitialized ledger")

	// ErrNilReservation identifies a method call on a nil reservation receiver.
	ErrNilReservation = errors.New("capacity: nil reservation")

	// ErrInvalidReservation identifies a reservation not created by its ledger.
	ErrInvalidReservation = errors.New("capacity: invalid reservation")

	// ErrReservationReleased identifies a strict double-release attempt.
	ErrReservationReleased = errors.New("capacity: reservation already released")

	// ErrReservedUnderflow identifies impossible internal reserved-state underflow.
	ErrReservedUnderflow = errors.New("capacity: reserved amount underflow")
)

// Error is the structured diagnostic returned or panicked by capacity.
type Error struct {
	// Record stores path, sentinel, reason, detail, and nested cause data.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable capacity diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Record.Format("capacity")
}

// Unwrap returns the broad sentinel error and optional nested cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Record.Unwrap()
}

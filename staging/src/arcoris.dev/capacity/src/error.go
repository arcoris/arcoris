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
	"fmt"
)

// Sentinel errors identify broad capacity failure classes.
var (
	// ErrInvalidResource identifies an invalid resource identity.
	ErrInvalidResource = errors.New("capacity: invalid resource")

	// ErrZeroAmount identifies a vector or demand entry with zero amount.
	ErrZeroAmount = errors.New("capacity: zero amount")

	// ErrDuplicateResource identifies repeated resources in constructor input.
	ErrDuplicateResource = errors.New("capacity: duplicate resource")

	// ErrEmptyDemand identifies an empty reservation demand.
	ErrEmptyDemand = errors.New("capacity: empty demand")

	// ErrInvalidVector identifies a vector that is not canonical.
	ErrInvalidVector = errors.New("capacity: invalid vector")

	// ErrInvalidDemand identifies a demand that is empty or not canonical.
	ErrInvalidDemand = errors.New("capacity: invalid demand")

	// ErrInvalidVectorState identifies vector source state that is not canonical.
	ErrInvalidVectorState = errors.New("capacity: invalid vector state")

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

// Error is a compact structured capacity error.
//
// Error wraps one sentinel so callers can use errors.Is while still receiving a
// stable path and human-facing detail in diagnostics. It intentionally remains
// package-local in shape; capacity does not own a generic diagnostics framework.
type Error struct {
	// Path identifies the input field or owner state that failed validation.
	Path string

	// Err is the sentinel error matched by errors.Is.
	Err error

	// Detail describes the failure without carrying runtime request data.
	Detail string

	// Cause optionally carries a lower-level validation failure.
	Cause error
}

// Error returns a compact human-readable capacity diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	message := e.Err.Error()
	if e.Path != "" {
		message = fmt.Sprintf("%s at %s", message, e.Path)
	}
	if e.Detail != "" {
		message = fmt.Sprintf("%s: %s", message, e.Detail)
	}
	if e.Cause != nil {
		message = fmt.Sprintf("%s: %v", message, e.Cause)
	}

	return message
}

// Unwrap returns the sentinel and optional nested cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	if e.Cause == nil {
		return e.Err
	}

	return errors.Join(e.Err, e.Cause)
}

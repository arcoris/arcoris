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

// ErrorReason refines broad capacity sentinels with stable diagnostics.
type ErrorReason string

// Error reasons are stable machine-readable refinements for capacity errors.
const (
	// ErrorReasonInvalidResource reports a resource identifier grammar failure.
	ErrorReasonInvalidResource ErrorReason = "invalid_resource"

	// ErrorReasonZeroAmount reports an explicit zero where a positive amount is required.
	ErrorReasonZeroAmount ErrorReason = "zero_amount"

	// ErrorReasonDuplicateResource reports repeated entries for one resource.
	ErrorReasonDuplicateResource ErrorReason = "duplicate_resource"

	// ErrorReasonEmptyDemand reports an absent reservation demand.
	ErrorReasonEmptyDemand ErrorReason = "empty_demand"

	// ErrorReasonAmountOverflow reports checked unsigned addition overflow.
	ErrorReasonAmountOverflow ErrorReason = "amount_overflow"

	// ErrorReasonAmountUnderflow reports checked unsigned subtraction underflow.
	ErrorReasonAmountUnderflow ErrorReason = "amount_underflow"

	// ErrorReasonInvalidVector reports a non-canonical resource vector.
	ErrorReasonInvalidVector ErrorReason = "invalid_vector"

	// ErrorReasonInvalidDemand reports an invalid demand value.
	ErrorReasonInvalidDemand ErrorReason = "invalid_demand"

	// ErrorReasonInvalidState reports invalid source accounting state.
	ErrorReasonInvalidState ErrorReason = "invalid_state"

	// ErrorReasonNilLedger reports a nil ledger receiver.
	ErrorReasonNilLedger ErrorReason = "nil_ledger"

	// ErrorReasonUninitializedLedger reports a zero-value stateful ledger owner.
	ErrorReasonUninitializedLedger ErrorReason = "uninitialized_ledger"

	// ErrorReasonNilReservation reports a nil reservation receiver.
	ErrorReasonNilReservation ErrorReason = "nil_reservation"

	// ErrorReasonInvalidReservation reports a detached or zero reservation.
	ErrorReasonInvalidReservation ErrorReason = "invalid_reservation"

	// ErrorReasonReservationReleased reports a strict double-release attempt.
	ErrorReasonReservationReleased ErrorReason = "reservation_released"

	// ErrorReasonReservedUnderflow reports impossible internal reserved-state underflow.
	ErrorReasonReservedUnderflow ErrorReason = "reserved_underflow"
)

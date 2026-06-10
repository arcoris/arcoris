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

package value

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// API value construction sentinels classify broad construction failures.
//
// Constructors return structured *Error values that unwrap to these sentinels,
// so callers can use errors.Is without depending on human-readable details.
var (
	// ErrInvalidValue classifies malformed concrete API values.
	ErrInvalidValue = errors.New("invalid API value")
	// ErrInvalidRecord classifies malformed record payload values.
	ErrInvalidRecord = errors.New("invalid record value")
	// ErrInvalidList classifies malformed list payload values.
	ErrInvalidList = errors.New("invalid list value")
	// ErrInvalidRecordMember classifies malformed record member inputs.
	ErrInvalidRecordMember = errors.New("invalid record member")
	// ErrDuplicateMemberName classifies repeated record member names.
	ErrDuplicateMemberName = errors.New("duplicate record member name")
	// ErrEmptyMemberName classifies empty record member names.
	ErrEmptyMemberName = errors.New("empty record member name")
	// ErrInvalidFloat classifies NaN and infinity float inputs.
	ErrInvalidFloat = errors.New("invalid float value")
	// ErrInvalidDecimal classifies malformed decimal text inputs.
	ErrInvalidDecimal = errors.New("invalid decimal value")
	// ErrInvalidDate classifies impossible calendar date inputs.
	ErrInvalidDate = errors.New("invalid date value")
	// ErrInvalidTime classifies impossible time-of-day inputs.
	ErrInvalidTime = errors.New("invalid time-of-day value")
)

// Error is the structured diagnostic returned by value constructors.
//
// Err and Reason are stable programmatic classification members. Detail is for
// human diagnostics. Cause is reserved for nested construction failures when a
// future constructor wraps another construction path.
type Error struct {
	// Record stores the shared path, sentinel, reason, detail, and cause fields.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable value diagnostic.
//
// The string is intentionally diagnostic only. Callers should use errors.Is,
// errors.As, Path, and Reason for stable behavior.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("value")
}

// Unwrap preserves both broad value classification and nested causes.
//
// For example, a duplicate record member unwraps to ErrDuplicateMemberName,
// ErrInvalidRecord, and ErrInvalidValue. This keeps callers free to handle
// either specific failures or broad value-construction failures.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	base := e.Err
	if isInvalidValueFailure(e.Err) {
		base = errors.Join(ErrInvalidValue, e.Err)
	}
	if isInvalidRecordFailure(e.Err) {
		base = errors.Join(base, ErrInvalidRecord)
	}

	return diagnostic.JoinRecord[ErrorReason](base, e.Cause).Unwrap()
}

// isInvalidValueFailure reports whether err is one concrete value failure.
//
// These errors unwrap to ErrInvalidValue so callers can catch all construction
// failures with one sentinel.
func isInvalidValueFailure(err error) bool {
	switch err {
	case ErrInvalidRecord,
		ErrInvalidList,
		ErrInvalidRecordMember,
		ErrDuplicateMemberName,
		ErrEmptyMemberName,
		ErrInvalidFloat,
		ErrInvalidDecimal,
		ErrInvalidDate,
		ErrInvalidTime:
		return true
	default:
		return false
	}
}

// isInvalidRecordFailure reports whether err belongs to record construction.
//
// Record member-specific failures also unwrap to ErrInvalidRecord.
func isInvalidRecordFailure(err error) bool {
	switch err {
	case ErrInvalidRecordMember,
		ErrDuplicateMemberName,
		ErrEmptyMemberName:
		return true
	default:
		return false
	}
}

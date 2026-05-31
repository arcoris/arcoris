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
	"strings"
)

// API value construction sentinels classify broad construction failures.
//
// Constructors return structured *Error values that unwrap to these sentinels,
// so callers can use errors.Is without depending on human-readable details.
var (
	// ErrInvalidValue classifies malformed concrete API values.
	ErrInvalidValue = errors.New("invalid API value")
	// ErrInvalidObject classifies malformed object payload values.
	ErrInvalidObject = errors.New("invalid object value")
	// ErrInvalidList classifies malformed list payload values.
	ErrInvalidList = errors.New("invalid list value")
	// ErrInvalidMap classifies malformed map payload values.
	ErrInvalidMap = errors.New("invalid map value")
	// ErrInvalidField classifies malformed object field inputs.
	ErrInvalidField = errors.New("invalid object field")
	// ErrInvalidEntry classifies malformed map entry inputs.
	ErrInvalidEntry = errors.New("invalid map entry")
	// ErrDuplicateName classifies repeated object field names.
	ErrDuplicateName = errors.New("duplicate object field name")
	// ErrDuplicateKey classifies repeated map entry keys.
	ErrDuplicateKey = errors.New("duplicate map key")
	// ErrEmptyName classifies empty object field names.
	ErrEmptyName = errors.New("empty object field name")
	// ErrEmptyKey classifies empty map keys.
	ErrEmptyKey = errors.New("empty map key")
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
// Err and Reason are stable programmatic classification fields. Detail is for
// human diagnostics. Cause is reserved for nested construction failures when a
// future constructor wraps another construction path.
type Error struct {
	// Path identifies the value location that failed construction.
	Path string
	// Err is the broad sentinel used with errors.Is.
	Err error
	// Reason gives stable machine-readable detail within Err.
	Reason ErrorReason
	// Detail gives human-readable context for logs and diagnostics.
	Detail string
	// Cause preserves nested construction failures.
	Cause error
}

// Error returns a compact human-readable value diagnostic.
//
// The string is intentionally diagnostic only. Callers should use errors.Is,
// errors.As, Path, and Reason for stable behavior.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	parts := []string{"value"}

	if e.Path != "" {
		parts = append(parts, e.Path)
	}

	if e.Err != nil {
		parts = append(parts, e.Err.Error())
	}

	if e.Reason != "" {
		parts = append(parts, string(e.Reason))
	}

	if e.Detail != "" {
		parts = append(parts, e.Detail)
	}

	return strings.Join(parts, ": ")
}

// Unwrap preserves both broad value classification and nested causes.
//
// For example, a duplicate object field unwraps to ErrDuplicateName,
// ErrInvalidObject, and ErrInvalidValue. This keeps callers free to handle
// either specific failures or broad value-construction failures.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	base := e.Err
	if isInvalidValueFailure(e.Err) {
		base = errors.Join(ErrInvalidValue, e.Err)
	}
	if isInvalidObjectFailure(e.Err) {
		base = errors.Join(base, ErrInvalidObject)
	}
	if isInvalidMapFailure(e.Err) {
		base = errors.Join(base, ErrInvalidMap)
	}

	if base != nil && e.Cause != nil {
		return errors.Join(base, e.Cause)
	}

	return base
}

// isInvalidValueFailure reports whether err is one concrete value failure.
//
// These errors unwrap to ErrInvalidValue so callers can catch all construction
// failures with one sentinel.
func isInvalidValueFailure(err error) bool {
	switch err {
	case ErrInvalidObject,
		ErrInvalidList,
		ErrInvalidMap,
		ErrInvalidField,
		ErrInvalidEntry,
		ErrDuplicateName,
		ErrDuplicateKey,
		ErrEmptyName,
		ErrEmptyKey,
		ErrInvalidFloat,
		ErrInvalidDecimal,
		ErrInvalidDate,
		ErrInvalidTime:
		return true
	default:
		return false
	}
}

// isInvalidObjectFailure reports whether err belongs to object construction.
//
// Object field-specific failures also unwrap to ErrInvalidObject.
func isInvalidObjectFailure(err error) bool {
	switch err {
	case ErrInvalidField,
		ErrDuplicateName,
		ErrEmptyName:
		return true
	default:
		return false
	}
}

// isInvalidMapFailure reports whether err belongs to map construction.
//
// Map entry-specific failures also unwrap to ErrInvalidMap.
func isInvalidMapFailure(err error) bool {
	switch err {
	case ErrInvalidEntry,
		ErrDuplicateKey,
		ErrEmptyKey:
		return true
	default:
		return false
	}
}

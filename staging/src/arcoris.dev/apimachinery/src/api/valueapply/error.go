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

package valueapply

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidRequest classifies request shapes that cannot be processed by
	// valueapply itself, such as failed ownership-state transformations.
	ErrInvalidRequest = errors.New("invalid value apply request")

	// ErrInvalidPath classifies malformed semantic base paths.
	ErrInvalidPath = errors.New("invalid field path")

	// ErrInvalidOwner classifies malformed field ownership identities.
	ErrInvalidOwner = errors.New("invalid field owner")

	// ErrInvalidValue classifies live or applied payloads rejected by
	// descriptor-aware validation.
	ErrInvalidValue = errors.New("invalid value")

	// ErrFieldSetFailed classifies failures while extracting applied semantic
	// field paths.
	ErrFieldSetFailed = errors.New("value field set extraction failed")

	// ErrCompareFailed classifies failures while comparing live and applied
	// values.
	ErrCompareFailed = errors.New("value comparison failed")

	// ErrConflict classifies changed applied fields that overlap ownership held
	// by other owners.
	ErrConflict = errors.New("field ownership conflict")

	// ErrUnsupportedTakeover classifies forced ownership takeovers that cannot
	// be represented precisely by fieldownership.State transformations.
	ErrUnsupportedTakeover = errors.New("unsupported ownership takeover")

	// ErrMergeFailed classifies failures while merging selected fields.
	ErrMergeFailed = errors.New("value merge failed")
)

// Error is the structured diagnostic returned for one apply failure.
type Error struct {
	// Record stores the common API diagnostic fields while keeping ErrorReason
	// package-local and strongly typed.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable valueapply diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("valueapply")
}

// Unwrap exposes the broad sentinel and nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

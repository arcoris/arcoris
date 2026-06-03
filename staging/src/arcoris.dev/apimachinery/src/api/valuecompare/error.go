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

package valuecompare

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

var (
	// ErrInvalidValue classifies concrete values that cannot be traversed.
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidDescriptor classifies malformed descriptors encountered defensively.
	ErrInvalidDescriptor = errors.New("invalid descriptor")

	// ErrKindMismatch classifies value kind / descriptor type mismatches.
	ErrKindMismatch = errors.New("value kind mismatch")

	// ErrUnknownField classifies actual object members rejected by descriptor policy.
	ErrUnknownField = errors.New("unknown field")

	// ErrUnresolvedRef classifies TypeRef descriptors that cannot resolve.
	ErrUnresolvedRef = errors.New("unresolved type reference")

	// ErrReferenceCycle classifies recursive or too-deep TypeRef traversal.
	ErrReferenceCycle = errors.New("type reference cycle")

	// ErrInvalidListKey classifies ListMap selector extraction failures.
	ErrInvalidListKey = errors.New("invalid list map key")

	// ErrDuplicateListKey classifies repeated ListMap selector identities.
	ErrDuplicateListKey = errors.New("duplicate list map key")
)

// Error is the structured diagnostic returned for one comparison failure.
type Error struct {
	// Record stores path, sentinel, reason, detail, and optional nested cause.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable value comparison diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("valuecompare")
}

// Unwrap preserves the broad sentinel and any nested cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

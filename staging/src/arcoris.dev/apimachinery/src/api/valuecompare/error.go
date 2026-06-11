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
	// ErrInvalidValue classifies concrete payload values that block traversal.
	ErrInvalidValue = errors.New("invalid value")

	// ErrInvalidDescriptor classifies unusable descriptor shapes or views.
	ErrInvalidDescriptor = errors.New("invalid descriptor")

	// ErrInvalidPath classifies malformed semantic paths supplied by callers or
	// constructed during comparison.
	ErrInvalidPath = errors.New("invalid field path")

	// ErrKindMismatch classifies concrete value kind / descriptor kind mismatches.
	ErrKindMismatch = errors.New("value kind mismatch")

	// ErrUnknownField classifies record members rejected by an object descriptor.
	ErrUnknownField = errors.New("unknown field")

	// ErrUnresolvedRef classifies DescriptorRef descriptors that cannot be resolved.
	ErrUnresolvedRef = errors.New("unresolved descriptor reference")

	// ErrReferenceCycle classifies recursive or over-depth DescriptorRef traversal.
	ErrReferenceCycle = errors.New("descriptor reference cycle")

	// ErrInvalidListKey classifies ListMap keys that cannot form selectors.
	ErrInvalidListKey = errors.New("invalid list map key")

	// ErrDuplicateListKey classifies repeated ListMap selector identities.
	ErrDuplicateListKey = errors.New("duplicate list map key")

	// ErrInvalidResult classifies malformed comparison result values.
	ErrInvalidResult = errors.New("invalid comparison result")
)

// Error is the structured diagnostic returned for one fail-fast comparison blocker.
type Error struct {
	// Record stores path, sentinel, reason, detail, and optional nested cause.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable valuecompare diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("valuecompare")
}

// Unwrap exposes the broad sentinel and any nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

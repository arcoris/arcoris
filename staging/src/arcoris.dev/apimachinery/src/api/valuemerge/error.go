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

package valuemerge

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
	// constructed during merge.
	ErrInvalidPath = errors.New("invalid field path")

	// ErrKindMismatch classifies concrete value kind / descriptor mismatches.
	ErrKindMismatch = errors.New("value kind mismatch")

	// ErrUnknownField classifies object members rejected by UnknownReject.
	ErrUnknownField = errors.New("unknown field")

	// ErrUnresolvedRef classifies DescriptorRef descriptors that cannot be resolved.
	ErrUnresolvedRef = errors.New("unresolved descriptor reference")

	// ErrReferenceCycle classifies recursive or over-depth DescriptorRef traversal.
	ErrReferenceCycle = errors.New("descriptor reference cycle")

	// ErrInvalidListKey classifies ListMap keys that cannot form selectors.
	ErrInvalidListKey = errors.New("invalid list map key")

	// ErrDuplicateListKey classifies repeated ListMap selector identities.
	ErrDuplicateListKey = errors.New("duplicate list map key")

	// ErrUnsupportedMerge classifies selected paths that cannot be merged under
	// the descriptor semantics without changing package policy.
	ErrUnsupportedMerge = errors.New("unsupported merge")
)

// Error is the structured diagnostic returned for one fail-fast merge blocker.
type Error struct {
	// Record stores path, sentinel, reason, detail, and optional nested cause.
	diagnostic.Record[ErrorReason]
}

// Error returns a compact human-readable valuemerge diagnostic.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}

	return e.Record.Format("valuemerge")
}

// Unwrap exposes the broad sentinel and nested cause for errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Record.Unwrap()
}

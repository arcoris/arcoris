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
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/valuefieldset"
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

// ErrorReason gives stable machine-readable detail inside a broad error category.
type ErrorReason string

const (
	// ErrorReasonInvalidZero reports the uninitialized zero value.Value.
	ErrorReasonInvalidZero ErrorReason = "invalid_zero"

	// ErrorReasonInvalidDescriptor reports a malformed or unsupported descriptor.
	ErrorReasonInvalidDescriptor ErrorReason = "invalid_descriptor"

	// ErrorReasonInvalidPath reports a malformed semantic field path.
	ErrorReasonInvalidPath ErrorReason = "invalid_path"

	// ErrorReasonKindMismatch reports descriptor/value kind incompatibility.
	ErrorReasonKindMismatch ErrorReason = "kind_mismatch"

	// ErrorReasonUnknownField reports an undeclared object member rejected by policy.
	ErrorReasonUnknownField ErrorReason = "unknown_field"

	// ErrorReasonUnresolvedRef reports a DescriptorRef the resolver cannot load.
	ErrorReasonUnresolvedRef ErrorReason = "unresolved_ref"

	// ErrorReasonReferenceCycle reports recursive or over-depth DescriptorRef traversal.
	ErrorReasonReferenceCycle ErrorReason = "reference_cycle"

	// ErrorReasonMissingListKey reports a ListMap item missing a key field.
	ErrorReasonMissingListKey ErrorReason = "missing_list_key"

	// ErrorReasonInvalidListKey reports a ListMap key that cannot form a selector.
	ErrorReasonInvalidListKey ErrorReason = "invalid_list_key"

	// ErrorReasonDuplicateListKey reports repeated ListMap selector identity.
	ErrorReasonDuplicateListKey ErrorReason = "duplicate_list_key"
)

// errorAt creates a compare diagnostic at the canonical semantic payload path.
func errorAt(path fieldpath.Path, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path.String(), err, reason, detail),
	}
}

// errorfAt creates a compare diagnostic with formatted detail text.
func errorfAt(path fieldpath.Path, err error, reason ErrorReason, format string, args ...any) error {
	return errorAt(path, err, reason, fmt.Sprintf(format, args...))
}

// wrapAt attaches valuecompare classification while preserving a lower-level cause.
func wrapAt(
	path fieldpath.Path,
	err error,
	reason ErrorReason,
	detail string,
	cause error,
) error {
	return &Error{
		Record: diagnostic.WrapRecord(path.String(), err, reason, detail, cause),
	}
}

// compareFieldSetError translates valuefieldset subtree failures into this package's error model.
func compareFieldSetError(path fieldpath.Path, err error) error {
	var fieldSetError *valuefieldset.Error
	if !errors.As(err, &fieldSetError) {
		return wrapAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidZero,
			"field-set extraction failed",
			err,
		)
	}

	sentinel, reason := compareFieldSetErrorKind(fieldSetError.Record.Err, fieldSetError.Reason)
	return &Error{
		Record: diagnostic.WrapRecord(
			fieldSetError.Path,
			sentinel,
			reason,
			fieldSetError.Detail,
			fieldSetError.Cause,
		),
	}
}

// compareFieldSetErrorKind maps public field-set sentinels to compare sentinels.
func compareFieldSetErrorKind(err error, reason valuefieldset.ErrorReason) (error, ErrorReason) {
	switch {
	case errors.Is(err, valuefieldset.ErrInvalidValue):
		return ErrInvalidValue, ErrorReasonInvalidZero
	case errors.Is(err, valuefieldset.ErrInvalidDescriptor):
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case errors.Is(err, valuefieldset.ErrKindMismatch):
		return ErrKindMismatch, ErrorReasonKindMismatch
	case errors.Is(err, valuefieldset.ErrUnknownField):
		return ErrUnknownField, ErrorReasonUnknownField
	case errors.Is(err, valuefieldset.ErrUnresolvedRef):
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case errors.Is(err, valuefieldset.ErrReferenceCycle):
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	case errors.Is(err, valuefieldset.ErrDuplicateListKey):
		return ErrDuplicateListKey, ErrorReasonDuplicateListKey
	case errors.Is(err, valuefieldset.ErrInvalidListKey):
		return ErrInvalidListKey, compareFieldSetListKeyReason(reason)
	default:
		return ErrInvalidValue, ErrorReasonInvalidZero
	}
}

// compareFieldSetListKeyReason preserves the more specific missing-key reason.
func compareFieldSetListKeyReason(reason valuefieldset.ErrorReason) ErrorReason {
	if reason == valuefieldset.ErrorReasonMissingListKey {
		return ErrorReasonMissingListKey
	}

	return ErrorReasonInvalidListKey
}

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

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/valuefieldset"
)

// compareFieldSetError translates valuefieldset subtree failures into this package's error model.
func compareFieldSetError(path fieldpath.Path, err error) error {
	var fieldSetError *valuefieldset.Error
	if !errors.As(err, &fieldSetError) {
		return wrapAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidValueKind,
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
	case errors.Is(err, valuefieldset.ErrInvalidPath):
		return ErrInvalidPath, ErrorReasonInvalidPath
	case errors.Is(err, valuefieldset.ErrInvalidValue):
		return compareFieldSetValueErrorReason(reason)
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
		return ErrInvalidValue, ErrorReasonInvalidValueKind
	}
}

// compareFieldSetValueErrorReason preserves valuefieldset's precise path errors.
func compareFieldSetValueErrorReason(reason valuefieldset.ErrorReason) (error, ErrorReason) {
	switch reason {
	case valuefieldset.ErrorReasonInvalidFieldName:
		return ErrInvalidPath, ErrorReasonInvalidFieldName
	case valuefieldset.ErrorReasonInvalidMapKey:
		return ErrInvalidPath, ErrorReasonInvalidMapKey
	case valuefieldset.ErrorReasonInvalidZero:
		return ErrInvalidValue, ErrorReasonInvalidZero
	default:
		return ErrInvalidValue, ErrorReasonInvalidValueKind
	}
}

// compareFieldSetListKeyReason preserves the more specific missing-key reason.
func compareFieldSetListKeyReason(reason valuefieldset.ErrorReason) ErrorReason {
	if reason == valuefieldset.ErrorReasonMissingListKey {
		return ErrorReasonMissingListKey
	}

	return ErrorReasonInvalidListKey
}

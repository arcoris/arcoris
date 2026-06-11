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

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/valuecompare"
	"arcoris.dev/apimachinery/api/valuefieldset"
	"arcoris.dev/apimachinery/api/valuemerge"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

// wrapValidationError classifies valuevalidation failures while preserving the
// original validation error as cause.
func wrapValidationError(path fieldpath.Path, detail string, err error) error {
	sentinel, reason := classifyValidationError(err)
	return wrapAt(path, sentinel, reason, detail, err)
}

// wrapFieldSetError classifies valuefieldset failures while preserving field
// extraction stage context.
func wrapFieldSetError(path fieldpath.Path, err error) error {
	sentinel, reason := classifyFieldSetError(err)
	return wrapAt(path, sentinel, reason, "applied field extraction failed", err)
}

// wrapCompareError classifies valuecompare failures while preserving comparison
// stage context.
func wrapCompareError(path fieldpath.Path, err error) error {
	sentinel, reason := classifyCompareError(err)
	return wrapAt(path, sentinel, reason, "value comparison failed", err)
}

// wrapMergeError classifies valuemerge failures while preserving merge stage context.
func wrapMergeError(path fieldpath.Path, err error) error {
	sentinel, reason := classifyMergeError(err)
	return wrapAt(path, sentinel, reason, "value merge failed", err)
}

// classifyValidationError maps valuevalidation diagnostics to valueapply policy
// categories without discarding the original cause.
func classifyValidationError(err error) (error, ErrorReason) {
	switch {
	case errors.Is(err, valuevalidation.ErrInvalidPath):
		return ErrInvalidPath, ErrorReasonInvalidPath
	case errors.Is(err, valuevalidation.ErrInvalidDescriptor):
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case errors.Is(err, valuevalidation.ErrUnresolvedRef):
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case errors.Is(err, valuevalidation.ErrReferenceCycle):
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	default:
		return ErrInvalidValue, ErrorReasonInvalidValue
	}
}

// classifyFieldSetError maps ownership-field extraction failures to valueapply
// policy categories.
func classifyFieldSetError(err error) (error, ErrorReason) {
	switch {
	case errors.Is(err, valuefieldset.ErrInvalidPath):
		return ErrInvalidPath, ErrorReasonInvalidPath
	case errors.Is(err, valuefieldset.ErrInvalidDescriptor):
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case errors.Is(err, valuefieldset.ErrUnresolvedRef):
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case errors.Is(err, valuefieldset.ErrReferenceCycle):
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	default:
		return ErrFieldSetFailed, ErrorReasonFieldSetFailed
	}
}

// classifyCompareError maps semantic comparison failures to valueapply policy
// categories.
func classifyCompareError(err error) (error, ErrorReason) {
	switch {
	case errors.Is(err, valuecompare.ErrInvalidPath):
		return ErrInvalidPath, ErrorReasonInvalidPath
	case errors.Is(err, valuecompare.ErrInvalidDescriptor):
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case errors.Is(err, valuecompare.ErrUnresolvedRef):
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case errors.Is(err, valuecompare.ErrReferenceCycle):
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	default:
		return ErrCompareFailed, ErrorReasonCompareFailed
	}
}

// classifyMergeError maps selected-field merge failures to valueapply policy
// categories.
func classifyMergeError(err error) (error, ErrorReason) {
	switch {
	case errors.Is(err, valuemerge.ErrInvalidPath):
		return ErrInvalidPath, ErrorReasonInvalidPath
	case errors.Is(err, valuemerge.ErrInvalidDescriptor):
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case errors.Is(err, valuemerge.ErrUnresolvedRef):
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case errors.Is(err, valuemerge.ErrReferenceCycle):
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	case errors.Is(err, valuemerge.ErrUnsupportedMerge):
		return ErrUnsupportedMerge, ErrorReasonUnsupportedMerge
	default:
		return ErrMergeFailed, ErrorReasonMergeFailed
	}
}

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

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/listmapkey"
)

// duplicateListMapEntryError reports repeated selector identity with both indexes.
func duplicateListMapEntryError(
	selectorPath fieldpath.Path,
	firstPath fieldpath.Path,
	duplicatePath fieldpath.Path,
) error {
	return errorfAt(
		selectorPath,
		ErrDuplicateListKey,
		ErrorReasonDuplicateListKey,
		"duplicate list map key; first occurrence at %s, duplicate at %s",
		firstPath,
		duplicatePath,
	)
}

// mergeListMapKeyError maps shared ListMap key failures to merge diagnostics.
func mergeListMapKeyError(path fieldpath.Path, err error) error {
	var keyError *listmapkey.Error
	if !errors.As(err, &keyError) {
		return wrapAt(
			path,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map key extraction failed",
			err,
		)
	}

	sentinel, reason := mergeListMapKeyErrorKind(keyError.Kind)
	if keyError.Cause != nil {
		return wrapAt(keyError.Path, sentinel, reason, keyError.Detail, keyError.Cause)
	}

	return errorAt(keyError.Path, sentinel, reason, keyError.Detail)
}

// mergeListMapKeyErrorKind maps selector failure kinds to merge reasons.
func mergeListMapKeyErrorKind(kind listmapkey.FailureKind) (error, ErrorReason) {
	switch kind {
	case listmapkey.FailureInvalidDescriptor:
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case listmapkey.FailureUnresolvedRef:
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case listmapkey.FailureReferenceCycle:
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	case listmapkey.FailureMissingKey:
		return ErrInvalidListKey, ErrorReasonMissingListKey
	case listmapkey.FailureNullKey,
		listmapkey.FailureKeyKindMismatch,
		listmapkey.FailureKeyIntegerRange,
		listmapkey.FailureInvalidSelector,
		listmapkey.FailureItemKindMismatch:
		return ErrInvalidListKey, ErrorReasonInvalidListKey
	default:
		return ErrInvalidListKey, ErrorReasonInvalidListKey
	}
}

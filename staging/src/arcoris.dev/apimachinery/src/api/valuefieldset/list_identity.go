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

package valuefieldset

import (
	"errors"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/listmapkey"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// listMapSelector extracts one stable ListMap selector for field-set extraction.
//
// Field-set extraction assumes the payload has already passed valuevalidation.
// It therefore cannot fall back to index paths when selector extraction fails:
// without a stable selector the resulting field set would not be semantic.
func (e *extractor) listMapSelector(
	indexPath fieldpath.Path,
	item value.Value,
	element types.Type,
	keys []types.FieldName,
) (fieldpath.Selector, error) {
	itemSelector, err := listmapkey.ExtractSelector(
		indexPath,
		item,
		element,
		keys,
		listmapkey.Options{
			Resolver: e.resolver,
			MaxDepth: e.maxDepth,
		},
	)
	if err != nil {
		return fieldpath.Selector{}, fieldSetListMapKeyError(err)
	}

	return itemSelector, nil
}

// fieldSetListMapKeyError maps shared ListMap key failures to public field-set errors.
func fieldSetListMapKeyError(err error) error {
	keyError, ok := listMapKeyFailure(err)
	if !ok {
		return err
	}

	sentinel, reason := fieldSetListMapKeyErrorKind(keyError.Kind)
	if keyError.Cause != nil {
		return wrapAt(
			keyError.Path,
			sentinel,
			reason,
			keyError.Detail,
			keyError.Cause,
		)
	}

	return errorAt(keyError.Path, sentinel, reason, keyError.Detail)
}

// listMapKeyFailure extracts the shared ListMap key diagnostic from err.
func listMapKeyFailure(err error) (*listmapkey.Error, bool) {
	var keyError *listmapkey.Error
	if errors.As(err, &keyError) {
		return keyError, true
	}

	return nil, false
}

// fieldSetListMapKeyErrorKind maps internal failure kinds to field-set diagnostics.
func fieldSetListMapKeyErrorKind(kind listmapkey.FailureKind) (error, ErrorReason) {
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

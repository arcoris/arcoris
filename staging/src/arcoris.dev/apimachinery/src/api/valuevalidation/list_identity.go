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

package valuevalidation

import (
	"errors"

	"arcoris.dev/apimachinery/api/internal/listmapkey"
)

// reportListMapKeyFailure maps shared ListMap key extraction failures into
// validation diagnostics.
//
// Ordinary payload failures return false so validateListMap can validate the
// item at its physical index path. That recursive pass reports key and nested
// field problems with the same descriptor-aware rules used everywhere else.
func (v *validator) reportListMapKeyFailure(err error) bool {
	if listmapkey.IsPayloadFailure(err) {
		return false
	}

	keyError, ok := listMapKeyFailure(err)
	if !ok {
		return false
	}

	validationError, reason := listMapKeyErrorKind(keyError.Kind)
	if keyError.Cause != nil {
		v.wrap(
			keyError.Path,
			validationError,
			reason,
			keyError.Detail,
			keyError.Cause,
		)
		return true
	}

	v.add(keyError.Path, validationError, reason, keyError.Detail)
	return true
}

// listMapKeyFailure extracts the shared ListMap key diagnostic from err.
func listMapKeyFailure(err error) (*listmapkey.Error, bool) {
	var keyError *listmapkey.Error
	if errors.As(err, &keyError) {
		return keyError, true
	}

	return nil, false
}

// listMapKeyErrorKind maps internal failure kinds to public validation errors.
func listMapKeyErrorKind(kind listmapkey.FailureKind) (error, ErrorReason) {
	switch kind {
	case listmapkey.FailureInvalidDescriptor:
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case listmapkey.FailureUnresolvedRef:
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case listmapkey.FailureReferenceCycle:
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	case listmapkey.FailureInvalidSelector:
		return ErrInvalidListKey, ErrorReasonInvalidListKey
	default:
		return ErrInvalidListKey, ErrorReasonInvalidListKey
	}
}

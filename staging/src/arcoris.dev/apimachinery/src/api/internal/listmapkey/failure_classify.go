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

package listmapkey

import "errors"

// IsPayloadFailure reports whether err describes ordinary concrete payload data
// that prevented selector extraction.
//
// valuevalidation can usually ignore these key extraction diagnostics and
// validate the same item at its physical list index to produce better
// field-level messages.
func IsPayloadFailure(err error) bool {
	extractionError, ok := asError(err)
	if !ok {
		return false
	}

	switch extractionError.Kind {
	case FailureItemKindMismatch,
		FailureMissingKey,
		FailureNullKey,
		FailureKeyKindMismatch,
		FailureKeyIntegerRange:
		return true
	default:
		return false
	}
}

// IsDescriptorFailure reports whether err describes descriptor, resolver, or
// selector construction state rather than ordinary payload data.
//
// Callers should surface these failures directly because falling back to an
// index path would hide a problem in descriptor preparation or identity rules.
func IsDescriptorFailure(err error) bool {
	extractionError, ok := asError(err)
	if !ok {
		return false
	}

	switch extractionError.Kind {
	case FailureInvalidDescriptor,
		FailureUnresolvedRef,
		FailureReferenceCycle,
		FailureInvalidSelector:
		return true
	default:
		return false
	}
}

// asError extracts this package's structured error from an arbitrary error.
func asError(err error) (*Error, bool) {
	var extractionError *Error
	if errors.As(err, &extractionError) {
		return extractionError, true
	}

	return nil, false
}

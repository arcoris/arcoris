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

package objectstorewatch

import "errors"

// invalidSnapshotError builds a snapshot diagnostic at path.
//
// The returned error always matches ErrInvalidSnapshot and preserves cause, so
// callers can distinguish the bridge-layer failure from the lower objectstore
// validation reason without losing either one.
func invalidSnapshotError(path string, cause error) error {
	return objectStoreWatchError(
		path,
		ErrorReasonInvalidSnapshot,
		[]error{ErrInvalidSnapshot},
		cause,
	)
}

// invalidBoundaryError builds a boundary diagnostic at path.
//
// Boundary errors are intentionally separate from snapshot errors because a
// malformed collection/revision boundary may be produced without a list result.
func invalidBoundaryError(path string, cause error) error {
	return objectStoreWatchError(
		path,
		ErrorReasonInvalidBoundary,
		[]error{ErrInvalidBoundary},
		cause,
	)
}

// objectStoreWatchError builds the package's standard multi-layer diagnostic.
//
// The first layer is one or more broad sentinels for errors.Is. The second is a
// structured *Error for errors.As. The remaining layers are lower causes, such
// as objectstore validation failures. Keeping all layers joined makes every
// caller-facing error both machine-checkable and human-readable.
func objectStoreWatchError(path string, reason ErrorReason, sentinels []error, causes ...error) error {
	cause := joinedCauses(causes...)
	if cause == nil && len(sentinels) > 0 {
		cause = sentinels[len(sentinels)-1]
	}
	if cause == nil {
		cause = errors.New("object store watch contract violation")
	}

	diagnostic := &Error{Path: path, Reason: reason, Cause: cause}
	joined := make([]error, 0, len(sentinels)+len(causes)+1)
	joined = append(joined, sentinels...)
	joined = append(joined, diagnostic)
	for _, item := range causes {
		if item != nil {
			joined = append(joined, item)
		}
	}

	return errors.Join(joined...)
}

// joinedCauses combines non-nil causes while preserving each for errors.Is.
func joinedCauses(causes ...error) error {
	var joined error
	for _, item := range causes {
		if item == nil {
			continue
		}
		if joined == nil {
			joined = item
			continue
		}
		joined = errors.Join(joined, item)
	}

	return joined
}

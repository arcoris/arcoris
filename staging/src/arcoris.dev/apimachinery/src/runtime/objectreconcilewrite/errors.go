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

package objectreconcilewrite

import "errors"

var (
	// ErrInvalidCurrent reports a zero or malformed Current used to build a
	// lifecycle write request.
	ErrInvalidCurrent = errors.New("objectreconcilewrite: invalid current object")

	// ErrMissingRevision reports a current object state without the committed
	// object revision needed for optimistic concurrency.
	ErrMissingRevision = errors.New("objectreconcilewrite: missing object revision")

	// ErrInvalidRequest reports a malformed objectreconciler.Request.
	ErrInvalidRequest = errors.New("objectreconcilewrite: invalid request")

	// ErrInvalidSnapshot reports a snapshot read that cannot serve the requested
	// object key.
	ErrInvalidSnapshot = errors.New("objectreconcilewrite: invalid snapshot")
)

// errorWith preserves a package-level category while keeping lower-level
// validation details visible to errors.Is callers.
func errorWith(category error, cause error) error {
	if cause == nil {
		return category
	}

	return errors.Join(category, cause)
}

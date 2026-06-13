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

package objectstore

import "errors"

var (
	// ErrNotFound classifies a missing or tombstoned live object.
	ErrNotFound = errors.New("objectstore: not found")

	// ErrAlreadyExists classifies a create conflict with an existing live object.
	ErrAlreadyExists = errors.New("objectstore: already exists")

	// ErrConflict classifies a concurrent state transition conflict.
	ErrConflict = errors.New("objectstore: conflict")

	// ErrStaleRevision classifies an expected revision mismatch.
	ErrStaleRevision = errors.New("objectstore: stale revision")

	// ErrInvalidKey classifies malformed object store keys.
	ErrInvalidKey = errors.New("objectstore: invalid key")

	// ErrInvalidListRequest classifies malformed object store list requests.
	ErrInvalidListRequest = errors.New("objectstore: invalid list request")

	// ErrInvalidState classifies malformed object store state values.
	ErrInvalidState = errors.New("objectstore: invalid state")

	// ErrInvalidRevision classifies invalid or forged revisions.
	ErrInvalidRevision = errors.New("objectstore: invalid revision")

	// ErrNilContext classifies a nil context passed to a store operation.
	ErrNilContext = errors.New("objectstore: nil context")

	// ErrUninitializedStore classifies use of a nil or zero implementation.
	ErrUninitializedStore = errors.New("objectstore: uninitialized store")
)

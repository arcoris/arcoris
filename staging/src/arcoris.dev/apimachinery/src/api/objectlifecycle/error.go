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

package objectlifecycle

import (
	"errors"
)

var (
	// ErrInvalidRequest classifies malformed lifecycle operation input.
	ErrInvalidRequest = errors.New("objectlifecycle: invalid request")
	// ErrInvalidExecutor classifies missing executor dependencies.
	ErrInvalidExecutor = errors.New("objectlifecycle: invalid executor")
	// ErrResourceNotFound classifies resource resolver misses.
	ErrResourceNotFound = errors.New("objectlifecycle: resource not found")
	// ErrValidationFailed classifies objectvalidation failures.
	ErrValidationFailed = errors.New("objectlifecycle: validation failed")
	// ErrApplyFailed classifies non-conflict objectapply or ownership failures.
	ErrApplyFailed = errors.New("objectlifecycle: apply failed")
	// ErrConflict classifies field ownership or store transition conflicts.
	ErrConflict = errors.New("objectlifecycle: conflict")
	// ErrNotFound classifies missing live objects.
	ErrNotFound = errors.New("objectlifecycle: not found")
	// ErrAlreadyExists classifies create attempts for existing live objects.
	ErrAlreadyExists = errors.New("objectlifecycle: already exists")
	// ErrStaleRevision classifies stale expected revisions.
	ErrStaleRevision = errors.New("objectlifecycle: stale revision")
	// ErrStoreFailed classifies unexpected objectstore failures.
	ErrStoreFailed = errors.New("objectlifecycle: store failed")

	// ErrNilOption classifies nil constructor options.
	ErrNilOption = errors.New("objectlifecycle: nil option")
	// ErrNilStore classifies a missing objectstore.Store dependency.
	ErrNilStore = errors.New("objectlifecycle: nil store")
	// ErrNilResourceResolver classifies a missing resource resolver dependency.
	ErrNilResourceResolver = errors.New("objectlifecycle: nil resource resolver")
	// ErrNilDesiredValidator classifies a missing Desired surface validator.
	ErrNilDesiredValidator = errors.New("objectlifecycle: nil desired validator")
	// ErrNilContext classifies nil operation contexts.
	ErrNilContext = errors.New("objectlifecycle: nil context")
)

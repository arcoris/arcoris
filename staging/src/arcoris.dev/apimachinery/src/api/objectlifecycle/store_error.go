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

	"arcoris.dev/apimachinery/api/objectstore"
)

// mapStoreError converts objectstore concurrency errors into lifecycle errors.
func mapStoreError(op Operation, key objectstore.Key, err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, objectstore.ErrNotFound):
		return errorFor(op, ErrorReasonNotFound, key, ErrNotFound, err)
	case errors.Is(err, objectstore.ErrAlreadyExists):
		return errorFor(op, ErrorReasonAlreadyExists, key, ErrAlreadyExists, err)
	case errors.Is(err, objectstore.ErrStaleRevision):
		return errorFor(op, ErrorReasonStaleRevision, key, ErrStaleRevision, err)
	case errors.Is(err, objectstore.ErrConflict):
		return errorFor(op, ErrorReasonConflict, key, ErrConflict, err)
	case errors.Is(err, objectstore.ErrInvalidKey):
		return errorFor(op, ErrorReasonInvalidRequest, key, ErrInvalidRequest, err)
	case errors.Is(err, objectstore.ErrInvalidListRequest):
		return errorFor(op, ErrorReasonInvalidRequest, key, ErrInvalidRequest, err)
	case errors.Is(err, objectstore.ErrInvalidRevision):
		return errorFor(op, ErrorReasonInvalidExpectedRevision, key, ErrInvalidRequest, err)
	case errors.Is(err, objectstore.ErrInvalidState):
		return errorFor(op, ErrorReasonStoreInvalidState, key, ErrStoreFailed, err)
	case errors.Is(err, objectstore.ErrNilContext):
		return errorFor(op, ErrorReasonInvalidContext, key, ErrInvalidRequest, err)
	case errors.Is(err, objectstore.ErrUninitializedStore):
		return errorFor(op, ErrorReasonInvalidExecutor, key, ErrInvalidExecutor, err)
	default:
		return errorFor(op, ErrorReasonStoreFailed, key, ErrStoreFailed, err)
	}
}

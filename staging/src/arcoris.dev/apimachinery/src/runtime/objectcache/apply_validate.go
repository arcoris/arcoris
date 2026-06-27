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

package objectcache

import (
	"fmt"

	"arcoris.dev/apimachinery/api/objectstore"
)

// validateApply checks cache-local preconditions without mutating collection
// state. The caller must already have validated the objectstore.Change shape,
// collection membership, and revision monotonicity.
func (col collection) validateApply(change objectstore.Change) error {
	switch change.Kind {
	case objectstore.ChangeCreated:
		return col.validateCreate(change)
	case objectstore.ChangeUpdated:
		return col.validateUpdate(change)
	case objectstore.ChangeDeleted:
		return col.validateDelete(change)
	default:
		return errorWith(ErrInvalidChange, objectstore.ErrInvalidChange)
	}
}

// validateCreate rejects create changes for keys already materialized.
//
// A created change for an existing key is a source/cache sequencing problem,
// not a structural objectstore.Change error.
func (col collection) validateCreate(change objectstore.Change) error {
	if _, exists := col.items[change.Key]; exists {
		return invalidChangeStateError("create key already exists", change.Key)
	}

	return nil
}

// validateUpdate requires an existing key whose cached state exactly matches
// the change's Before side.
//
// Matching the full Before state prevents the cache from applying an update on
// top of a different observed object version with the same key.
func (col collection) validateUpdate(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("update key is missing", change.Key)
	}
	if !sameState(current.State, change.Before) {
		return invalidChangeStateError("update before state does not match cache", change.Key)
	}

	return nil
}

// validateDelete requires an existing key whose cached state exactly matches
// the change's tombstoned Before side.
//
// The delete revision is carried by change.Revision. Before must still match
// current live state so the cache does not remove the wrong object version.
func (col collection) validateDelete(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("delete key is missing", change.Key)
	}
	if !sameState(current.State, change.Before) {
		return invalidChangeStateError("delete before state does not match cache", change.Key)
	}

	return nil
}

// invalidChangeStateError reports a structurally valid change that cannot be
// applied to this cache's current materialized state.
func invalidChangeStateError(message string, key objectstore.Key) error {
	return errorWith(ErrInvalidChange, fmt.Errorf("%s: %s", message, key.String()))
}

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

// validateCreate rejects create changes for keys that are already materialized.
func (col collection) validateCreate(change objectstore.Change) error {
	if _, exists := col.items[change.Key]; exists {
		return invalidChangeStateError("create key already exists", change.Key)
	}

	return nil
}

// validateUpdate requires an existing key whose cached revision matches the
// change's Before state.
func (col collection) validateUpdate(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("update key is missing", change.Key)
	}
	if current.State.Revision != change.Before.Revision {
		return invalidChangeStateError("update before revision does not match cache", change.Key)
	}

	return nil
}

// validateDelete requires an existing key whose cached revision matches the
// live state tombstoned by the change.
func (col collection) validateDelete(change objectstore.Change) error {
	current, exists := col.items[change.Key]
	if !exists {
		return invalidChangeStateError("delete key is missing", change.Key)
	}
	if current.State.Revision != change.Before.Revision {
		return invalidChangeStateError("delete before revision does not match cache", change.Key)
	}

	return nil
}

// invalidChangeStateError reports a structurally valid change that does not
// match the cache's current materialized state.
func invalidChangeStateError(message string, key objectstore.Key) error {
	return fmt.Errorf("%w: %s: %s", ErrInvalidChange, message, key.String())
}

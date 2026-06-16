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

import "arcoris.dev/apimachinery/api/objectstore"

// validateApply checks cache-local preconditions without mutating collection
// state. The caller must already have validated the objectstore.Change shape and
// stale revision ordering.
func (col collection) validateApply(change objectstore.Change) error {
	switch change.Kind {
	case objectstore.ChangeCreated:
		return col.validateCreate(change)
	case objectstore.ChangeUpdated:
		return col.validateUpdate(change)
	case objectstore.ChangeDeleted:
		return col.validateDelete(change)
	default:
		return invalidChangeError(nil)
	}
}

// applyValidated dispatches a change that already passed validateApply.
func (col *collection) applyValidated(change objectstore.Change) {
	switch change.Kind {
	case objectstore.ChangeCreated:
		col.applyCreateValidated(change)
	case objectstore.ChangeUpdated:
		col.applyUpdateValidated(change)
	case objectstore.ChangeDeleted:
		col.applyDeleteValidated(change)
	}
}

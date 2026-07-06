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

import "arcoris.dev/apimachinery/api/objectstore"

// Current is the committed object state used to build lifecycle write requests.
//
// Current stores the object key and a detached objectstore.State. Its Revision
// is the state revision used as the Expected value for optimistic concurrency.
type Current struct {
	key   objectstore.Key
	state objectstore.State
}

// FromItem validates item and returns a Current built from its committed state.
func FromItem(item objectstore.ListItem) (Current, error) {
	if err := objectstore.ValidateListItem(item); err != nil {
		if item.State.Revision.IsZero() {
			return Current{}, errorWith(ErrMissingRevision, err)
		}

		return Current{}, err
	}
	if item.State.Revision.IsZero() {
		return Current{}, ErrMissingRevision
	}

	return Current{
		key:   item.Key,
		state: item.State.Clone(),
	}, nil
}

// Key returns the object identity used by lifecycle write requests.
func (c Current) Key() objectstore.Key {
	return c.key
}

// State returns a detached copy of the current committed state.
func (c Current) State() objectstore.State {
	return c.state.Clone()
}

// Revision returns the current object's committed state revision.
func (c Current) Revision() objectstore.Revision {
	return c.state.Revision
}

func (c Current) validate() error {
	if c.state.Revision.IsZero() {
		return errorWith(ErrInvalidCurrent, ErrMissingRevision)
	}
	if err := objectstore.ValidateListItem(objectstore.ListItem{
		Key:   c.key,
		State: c.state,
	}); err != nil {
		return errorWith(ErrInvalidCurrent, err)
	}

	return nil
}

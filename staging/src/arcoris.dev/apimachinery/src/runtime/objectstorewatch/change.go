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

import (
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// publishValidated stores change in retained history and fans it out live.
//
// The caller must hold Store.mu and must pass a structurally valid committed
// change. Matching streams receive cloned EventChanged values. Slow streams are
// terminated with continuity loss rather than causing the writer to block.
func (s *Store) publishValidated(change objectstore.Change) {
	change = change.Clone()
	s.history.append(change)
	for w := range s.watchers {
		if !objectstore.ChangeMatchesListRequest(change, w.request.Collection) {
			continue
		}
		if w.enqueueChange(change) {
			continue
		}
		w.terminate(streamOverflowError())
		delete(s.watchers, w)
	}
}

// streamOverflowError returns the terminal continuity error for slow streams.
func streamOverflowError() error {
	return continuityError(ErrStreamOverflow)
}

// changedEvent converts a committed change into a validated watch event.
//
// The objectwatch constructor clones and validates the payload, keeping stream
// delivery tied to objectwatch's event contract rather than duplicating it here.
func changedEvent(change objectstore.Change) (objectwatch.Event, error) {
	return objectwatch.Changed(change.Clone())
}

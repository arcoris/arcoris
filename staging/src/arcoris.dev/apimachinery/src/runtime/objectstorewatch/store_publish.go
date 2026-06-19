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

import "arcoris.dev/apimachinery/api/objectstore"

// publishLocked stores change in retained history and fans it out live.
//
// Store.mu must be held. The caller must pass a structurally valid committed
// change. Matching streams receive cloned EventChanged values. Slow streams are
// terminated with continuity loss rather than causing the writer to block.
func (s *Store) publishLocked(change objectstore.Change) error {
	change = change.Clone()
	s.history.append(change)

	event, err := changedEvent(change)
	if err != nil {
		return s.loseCommittedContinuityLocked(change.Revision, err)
	}

	for id, stream := range s.streams {
		if !objectstore.ChangeMatchesListRequest(change, stream.request.Collection) {
			continue
		}
		if stream.enqueue(event) {
			continue
		}
		err := streamOverflowError()
		s.unregisterLocked(id)
		stream.finish(err)
	}

	return nil
}

// streamOverflowError returns the terminal continuity error for slow streams.
func streamOverflowError() error {
	return continuityError(ErrStreamOverflow)
}

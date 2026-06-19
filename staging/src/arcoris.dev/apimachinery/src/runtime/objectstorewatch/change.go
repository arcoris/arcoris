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

// createdChange builds the committed create transition captured by the wrapper.
func createdChange(key objectstore.Key, after objectstore.State) (objectstore.Change, error) {
	return objectstore.NewCreatedChange(key, after)
}

// updatedChange builds the committed update transition captured by the wrapper.
func updatedChange(
	key objectstore.Key,
	before objectstore.State,
	after objectstore.State,
) (objectstore.Change, error) {
	return objectstore.NewUpdatedChange(key, before, after)
}

// deletedChange builds the committed delete transition captured by the wrapper.
//
// The delete revision must come from DeleteResult.Revision, not from the
// deleted live state, because the latter is the previous live revision.
func deletedChange(key objectstore.Key, result objectstore.DeleteResult) (objectstore.Change, error) {
	return objectstore.NewDeletedChange(key, result.Deleted, result.Revision)
}

// changedEvent converts a committed change into a validated watch event.
//
// The objectwatch constructor clones and validates the payload, keeping stream
// delivery tied to objectwatch's event contract rather than duplicating it here.
func changedEvent(change objectstore.Change) (objectwatch.Event, error) {
	return objectwatch.Changed(change.Clone())
}

// replayEvents converts retained changes into detached EventChanged replay data.
func replayEvents(changes []objectstore.Change) ([]objectwatch.Event, error) {
	if changes == nil {
		return nil, nil
	}

	events := make([]objectwatch.Event, len(changes))
	for i, change := range changes {
		event, err := changedEvent(change)
		if err != nil {
			return nil, err
		}
		events[i] = event
	}

	return events, nil
}

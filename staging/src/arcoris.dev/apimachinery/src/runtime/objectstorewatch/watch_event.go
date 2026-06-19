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

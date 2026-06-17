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

package objectwatch

import "arcoris.dev/apimachinery/api/objectstore"

// Changed constructs an event carrying one committed object transition.
//
// The change is cloned before storage in the Event so later caller mutation of
// the input Change cannot alter the event payload.
func Changed(change objectstore.Change) (Event, error) {
	event := Event{Kind: EventChanged, Revision: change.Revision, Change: change.Clone()}
	if err := event.Validate(); err != nil {
		return Event{}, err
	}

	return event, nil
}

// Bookmark constructs a progress-only event.
//
// A bookmark is useful only as a continuity/progress boundary. It never
// describes an object mutation and must never be applied to cached object
// state.
func Bookmark(revision objectstore.Revision) (Event, error) {
	event := Event{Kind: EventBookmark, Revision: revision}
	if err := event.Validate(); err != nil {
		return Event{}, err
	}

	return event, nil
}

// RestartRequired constructs a terminal continuity-loss event.
//
// A zero revision means the source cannot identify a useful last-known
// progress boundary. A non-zero revision names the last progress boundary the
// source is willing to report. Consumers must relist before trusting future
// state after this event.
func RestartRequired(reason RestartReason, revision objectstore.Revision) (Event, error) {
	event := Event{Kind: EventRestartRequired, Revision: revision, Restart: reason}
	if err := event.Validate(); err != nil {
		return Event{}, err
	}

	return event, nil
}

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

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestChangedConstructorsAcceptValidChanges(t *testing.T) {
	changes := []objectstore.Change{
		watchCreatedChange(1),
		watchUpdatedChange(1, 2),
		watchDeletedChange(1, 2),
	}

	for _, change := range changes {
		event, err := Changed(change)
		requireNoError(t, err)
		requireNoError(t, event.Validate())
		if event.Kind != EventChanged || event.Revision != change.Revision {
			t.Fatalf("event = %#v; want changed revision %s", event, change.Revision)
		}
	}
}

func TestChangedRejectsInvalidChange(t *testing.T) {
	_, err := Changed(objectstore.Change{})

	requireErrorIs(t, err, ErrInvalidEvent)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
	requireWatchError(t, err, ErrorReasonInvalidEvent, "watch.event.change")
}

func TestBookmarkValidate(t *testing.T) {
	event, err := Bookmark(10)
	requireNoError(t, err)

	if event.Kind != EventBookmark || event.Revision != 10 {
		t.Fatalf("event = %#v; want bookmark revision 10", event)
	}
	requireNoError(t, event.Validate())
}

func TestBookmarkRejectsZeroRevision(t *testing.T) {
	_, err := Bookmark(0)

	requireErrorIs(t, err, ErrInvalidEvent)
	requireWatchError(t, err, ErrorReasonInvalidEvent, "watch.event.bookmark")
}

func TestRestartRequiredValidate(t *testing.T) {
	zeroRevision, err := RestartRequired(RestartHistoryUnavailable, 0)
	requireNoError(t, err)
	if zeroRevision.Revision != 0 {
		t.Fatalf("zero restart revision = %s; want zero", zeroRevision.Revision)
	}

	withRevision, err := RestartRequired(RestartSourceReset, 20)
	requireNoError(t, err)
	if withRevision.Revision != 20 {
		t.Fatalf("restart revision = %s; want 20", withRevision.Revision)
	}
}

func TestRestartRequiredRejectsZeroReason(t *testing.T) {
	_, err := RestartRequired(0, 0)

	requireErrorIs(t, err, ErrInvalidEvent)
	requireErrorIs(t, err, ErrInvalidRestart)
	requireWatchError(t, err, ErrorReasonInvalidRestart, "watch.event.restart")
}

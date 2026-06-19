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

package objectreflector

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

func createdChange(t testing.TB, key objectstore.Key, revision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewCreatedChange(key, testState(key, revision, "created"))
	requireNoError(t, err)

	return change
}

func updatedChange(t testing.TB, key objectstore.Key, beforeRevision, afterRevision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewUpdatedChange(
		key,
		testState(key, beforeRevision, "before"),
		testState(key, afterRevision, "after"),
	)
	requireNoError(t, err)

	return change
}

func deletedChange(t testing.TB, key objectstore.Key, beforeRevision, deleteRevision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewDeletedChange(key, testState(key, beforeRevision, "deleted"), deleteRevision)
	requireNoError(t, err)

	return change
}

func changedEvent(t testing.TB, change objectstore.Change) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Changed(change)
	requireNoError(t, err)

	return event
}

func progressEvent(t testing.TB, revision objectstore.Revision) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Progress(revision)
	requireNoError(t, err)

	return event
}

func restartEvent(t testing.TB) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.RestartRequired(objectwatch.RestartContinuityLost, 0)
	requireNoError(t, err)

	return event
}

func TestProcessRestartRequiresRelist(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	err := reflector.processEvent(context.Background(), restartEvent(t))

	if !isRelistRequired(err) {
		t.Fatalf("isRelistRequired(%v) = false; want true", err)
	}
}

func TestProcessEventRejectsInvalidEvent(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))

	err := reflector.processEvent(context.Background(), objectwatch.Event{})

	requireErrorIs(t, err, ErrInvalidEvent)
}

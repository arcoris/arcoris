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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestRunInitialSyncBuildsWatchFromBoundary(t *testing.T) {
	read := testRead(t, 7)
	source := &fakeListerWatcher{
		listResponses:  []listResponse{{read: read}},
		watchResponses: []watchResponse{{stream: terminalStream(context.Canceled)}},
	}
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())
	requireErrorIs(t, err, context.Canceled)

	if sink.replaceCount() != 1 {
		t.Fatalf("replace count = %d; want 1", sink.replaceCount())
	}
	if sink.changeCount() != 0 {
		t.Fatalf("change count = %d; want 0", sink.changeCount())
	}
	requests := source.recordedWatchRequests()
	if len(requests) != 1 {
		t.Fatalf("watch requests = %d; want 1", len(requests))
	}
	if requests[0].Collection != testCollection() {
		t.Fatalf("watch collection = %#v; want %#v", requests[0].Collection, testCollection())
	}
	if requests[0].Start.Revision != read.Revision() {
		t.Fatalf("watch revision = %s; want %s", requests[0].Start.Revision, read.Revision())
	}
	if !requests[0].AllowProgress {
		t.Fatalf("AllowProgress = false; want true")
	}
}

func TestRunDoesNotWatchWhenReplaceFails(t *testing.T) {
	errReplace := errors.New("replace failed")
	source := &fakeListerWatcher{listResponses: []listResponse{{read: testRead(t, 1)}}}
	sink := newRecordingSink(1)
	sink.replaceErr = errReplace
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())

	requireErrorIs(t, err, errReplace)
	if len(source.recordedWatchRequests()) != 0 {
		t.Fatalf("Watch was called after Replace failure")
	}
}

func TestRunApplyChangeErrorIsFatal(t *testing.T) {
	errApply := errors.New("apply failed")
	source := &fakeListerWatcher{
		listResponses: []listResponse{{read: testRead(t, 1)}},
		watchResponses: []watchResponse{{stream: streamWithEvents(
			changedEvent(t, createdChange(t, testKey("system", 1), 2)),
		)}},
	}
	sink := newRecordingSink(1)
	sink.applyErr = errApply
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())

	requireErrorIs(t, err, errApply)
	if reflector.lastApplied != 1 {
		t.Fatalf("lastApplied = %s; want 1", reflector.lastApplied)
	}
}

func TestRunCycleRejectsInvalidCollectionRead(t *testing.T) {
	source := &fakeListerWatcher{
		listResponses: []listResponse{{read: testReadForCollection(t, namespaceCollection("other"), 1)}},
	}
	reflector := newTestReflector(t, source, newRecordingSink(1))

	err := reflector.runCycle(context.Background())

	requireErrorIs(t, err, ErrInvalidEvent)
}

func TestRunCycleAppliesChangesInRevisionOrder(t *testing.T) {
	key := testKey("system", 1)
	source := &fakeListerWatcher{
		listResponses: []listResponse{{read: testRead(t, 1, listItem(key, 1, "listed"))}},
		watchResponses: []watchResponse{{stream: streamWithEvents(
			changedEvent(t, updatedChange(t, key, 1, 2)),
			changedEvent(t, deletedChange(t, key, 2, 3)),
		)}},
	}
	sink := newRecordingSink(2)
	reflector := newTestReflector(t, source, sink)

	err := reflector.runCycle(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)

	changes := sink.recordedChanges()
	if len(changes) != 2 {
		t.Fatalf("changes = %d; want 2", len(changes))
	}
	if changes[0].Revision != 2 || changes[1].Revision != 3 {
		t.Fatalf("change revisions = %s, %s; want 2, 3", changes[0].Revision, changes[1].Revision)
	}
}

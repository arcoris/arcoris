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
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestWatchHistoricalReplay(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	updated := updateObject(t, store, key, created.Revision, "updated")
	deleted := deleteObject(t, store, key, updated.Revision)

	tests := []struct {
		name      string
		after     objectstore.Revision
		revisions []objectstore.Revision
	}{
		{name: "after zero", after: 0, revisions: []objectstore.Revision{created.Revision, updated.Revision, deleted.Revision}},
		{name: "after create", after: created.Revision, revisions: []objectstore.Revision{updated.Revision, deleted.Revision}},
		{name: "after update", after: updated.Revision, revisions: []objectstore.Revision{deleted.Revision}},
		{name: "after delete", after: deleted.Revision},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stream := watchAfter(t, store, testCollection(), tt.after)
			for _, revision := range tt.revisions {
				if event := nextEvent(t, stream); event.Revision != revision {
					t.Fatalf("revision = %s; want %s", event.Revision, revision)
				}
			}
			requireNoEvent(t, stream)
			requireNoError(t, stream.Close())
		})
	}
}

func TestWatchHistoricalReplayLargerThanStreamBuffer(t *testing.T) {
	store := testRuntimeStore(t, WithMaxHistory(10), WithStreamBuffer(1))
	first := createObject(t, store, testKey("system", 1), "one")
	second := createObject(t, store, testKey("system", 2), "two")
	third := createObject(t, store, testKey("system", 3), "three")

	stream := watchAfter(t, store, testCollection(), 0)

	for _, revision := range []objectstore.Revision{first.Revision, second.Revision, third.Revision} {
		if event := nextEvent(t, stream); event.Revision != revision {
			t.Fatalf("revision = %s; want %s", event.Revision, revision)
		}
	}
	requireNoEvent(t, stream)
}

func TestWatchHistoryCompaction(t *testing.T) {
	store := testRuntimeStore(t, WithMaxHistory(1))
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	updated := updateObject(t, store, key, created.Revision, "updated")

	start, err := objectwatch.AfterRevision(0)
	requireNoError(t, err)
	stream, err := store.Watch(context.Background(), objectwatch.Request{Collection: testCollection(), Start: start})
	if stream != nil {
		t.Fatalf("stream = %#v; want nil", stream)
	}
	requireErrorIs(t, err, objectwatch.ErrHistoryUnavailable)

	stream = watchAfter(t, store, testCollection(), created.Revision)
	event := nextEvent(t, stream)
	if event.Revision != updated.Revision {
		t.Fatalf("revision = %s; want %s", event.Revision, updated.Revision)
	}

	stream = watchAfter(t, store, testCollection(), updated.Revision)
	requireNoEvent(t, stream)
}

func TestListCollectionWatchContinuity(t *testing.T) {
	store := testRuntimeStore(t)
	read, err := store.ListCollection(context.Background(), testCollection())
	requireNoError(t, err)

	request, err := read.Boundary().WatchRequest(storewatchOptions())
	requireNoError(t, err)
	stream, err := store.Watch(context.Background(), request)
	requireNoError(t, err)

	created := createObject(t, store, testKey("system", 1), "created")
	event := nextEvent(t, stream)
	if event.Revision != created.Revision {
		t.Fatalf("revision = %s; want %s", event.Revision, created.Revision)
	}
}

func storewatchOptions() storewatchapi.WatchOptions {
	return storewatchapi.WatchOptions{}
}

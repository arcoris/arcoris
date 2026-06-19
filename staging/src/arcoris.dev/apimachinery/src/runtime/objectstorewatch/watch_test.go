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
	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestWatchRejectsInvalidRequest(t *testing.T) {
	store := testRuntimeStore(t)

	stream, err := store.Watch(context.Background(), objectwatch.Request{})

	if stream != nil {
		t.Fatalf("stream = %#v; want nil", stream)
	}
	requireErrorIs(t, err, objectwatch.ErrInvalidRequest)
}

func TestWatchReturnsStreamOnSuccess(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, testCollection(), 0)

	if stream == nil {
		t.Fatalf("stream is nil")
	}
	requireNoError(t, stream.Close())
}

func TestWatchStartAtCurrentDoesNotReplayHistory(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	stream := watchAtCurrent(t, store, testCollection())

	requireNoEvent(t, stream)
	updated := updateObject(t, store, key, created.Revision, "updated")
	event := nextEvent(t, stream)

	requireChangedEvent(t, event, objectwatch.Request{Collection: testCollection(), Start: objectwatch.AtCurrent()}, objectstore.ChangeUpdated, updated.Revision)
}

func TestWatchStartAtCurrentReturnsBackendListError(t *testing.T) {
	store := testRuntimeStore(t)

	stream, err := store.Watch(nil, objectwatch.Request{
		Collection: testCollection(),
		Start:      objectwatch.AtCurrent(),
	})

	if stream != nil {
		t.Fatalf("stream = %#v; want nil", stream)
	}
	requireErrorIs(t, err, objectstore.ErrNilContext)
}

func TestWatchStartAtCurrentRejectsInvalidListResult(t *testing.T) {
	store, err := New(invalidListResultBackend{Store: testBackend(t)})
	requireNoError(t, err)

	stream, err := store.Watch(context.Background(), objectwatch.Request{
		Collection: testCollection(),
		Start:      objectwatch.AtCurrent(),
	})

	if stream != nil {
		t.Fatalf("stream = %#v; want nil", stream)
	}
	requireErrorIs(t, err, objectstore.ErrInvalidListResult)
}

func TestWatchAllowsProgressButDoesNotEmitProgress(t *testing.T) {
	store := testRuntimeStore(t)
	start, err := objectwatch.AfterRevision(0)
	requireNoError(t, err)
	stream, err := store.Watch(context.Background(), objectwatch.Request{
		Collection:    testCollection(),
		Start:         start,
		AllowProgress: true,
	})
	requireNoError(t, err)

	requireNoEvent(t, stream)
}

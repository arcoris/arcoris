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
	"fmt"
	"sync"
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

func TestWatchFiltersByNamespaceAndResource(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, namespaceCollection("alpha"), 0)

	alpha := createObject(t, store, testKey("alpha", 1), "alpha")
	createObject(t, store, testKey("beta", 1), "beta")
	_, err := store.Create(context.Background(), otherResourceKey("alpha", 1), stateForKey(otherResourceKey("alpha", 1), "other"))
	requireNoError(t, err)

	event := nextEvent(t, stream)
	requireChangedEvent(t, event, watchRequestAfter(t, namespaceCollection("alpha"), 0), objectstore.ChangeCreated, alpha.Revision)
	requireNoEvent(t, stream)
}

func TestWatchAllNamespacesReceivesMultipleNamespaces(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, testCollection(), 0)

	first := createObject(t, store, testKey("alpha", 1), "alpha")
	second := createObject(t, store, testKey("beta", 1), "beta")

	if event := nextEvent(t, stream); event.Revision != first.Revision {
		t.Fatalf("first revision = %s; want %s", event.Revision, first.Revision)
	}
	if event := nextEvent(t, stream); event.Revision != second.Revision {
		t.Fatalf("second revision = %s; want %s", event.Revision, second.Revision)
	}
}

func TestConcurrentWatchAndCreate(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(16))
	request := watchRequestAfter(t, testCollection(), 0)
	errs := make(chan error, 16)
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			stream, err := store.Watch(context.Background(), request)
			if err != nil {
				errs <- fmt.Errorf("watch %d: %w", i, err)
				return
			}
			if err := stream.Close(); err != nil {
				errs <- fmt.Errorf("close %d: %w", i, err)
			}
		}(i)
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := testKey("system", i)
			_, err := store.Create(context.Background(), key, stateForKey(key, "created"))
			if err != nil {
				errs <- fmt.Errorf("create %d: %w", i, err)
			}
		}(i)
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		requireNoError(t, err)
	}
}

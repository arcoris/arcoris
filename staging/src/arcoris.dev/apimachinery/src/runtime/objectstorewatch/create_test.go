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

func TestCreatePublishesCreatedChange(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, testCollection(), 0)
	key := testKey("system", 1)

	created := createObject(t, store, key, "created")
	event := nextEvent(t, stream)

	requireChangedEvent(t, event, watchRequestAfter(t, testCollection(), 0), objectstore.ChangeCreated, created.Revision)
	if !event.Change.Key.Equal(key) {
		t.Fatalf("change key = %#v; want %#v", event.Change.Key, key)
	}
}

func TestCreatePublishesNothingWhenBackendFails(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	stream := watchAfter(t, store, testCollection(), created.Revision)

	_, err := store.Create(context.Background(), key, stateForKey(key, "duplicate"))
	if err == nil {
		t.Fatalf("Create() error = nil; want backend error")
	}

	requireNoEvent(t, stream)
}

func TestCreateContinuityLossTerminatesLiveStreams(t *testing.T) {
	backend := invalidCreateResultBackend{Store: testBackend(t)}
	store, err := New(backend)
	requireNoError(t, err)
	first := watchAfter(t, store, testCollection(), 0)
	second := watchAfter(t, store, testCollection(), 0)

	_, err = store.Create(context.Background(), testKey("system", 1), stateForKey(testKey("system", 1), "created"))
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)

	_, err = first.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	_, err = second.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
}

func TestCreateContinuityLossInvalidatesFutureHistoricalWatch(t *testing.T) {
	backend := invalidCreateStateBackend{Store: testBackend(t)}
	store, err := New(backend)
	requireNoError(t, err)
	stream := watchAfter(t, store, testCollection(), 0)
	key := testKey("system", 1)

	created, err := store.Create(context.Background(), key, stateForKey(key, "created"))
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
	_, err = stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)

	replay, err := store.Watch(context.Background(), watchRequestAfter(t, testCollection(), 0))
	if replay != nil {
		t.Fatalf("stream = %#v; want nil", replay)
	}
	requireErrorIs(t, err, objectwatch.ErrHistoryUnavailable)

	afterGap := watchAfter(t, store, testCollection(), created.Revision)
	updated := updateObject(t, store, key, created.Revision, "updated")
	event := nextEvent(t, afterGap)
	requireChangedEvent(t, event, watchRequestAfter(t, testCollection(), created.Revision), objectstore.ChangeUpdated, updated.Revision)
}

func TestCreateContinuityLossWithInvalidRevisionTaintsHistoricalWatch(t *testing.T) {
	backend := invalidCreateResultBackend{Store: testBackend(t)}
	store, err := New(backend)
	requireNoError(t, err)
	stream := watchAfter(t, store, testCollection(), 0)

	_, err = store.Create(context.Background(), testKey("system", 1), stateForKey(testKey("system", 1), "created"))
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
	_, err = stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)

	replay, err := store.Watch(context.Background(), watchRequestAfter(t, testCollection(), 0))
	if replay != nil {
		t.Fatalf("stream = %#v; want nil", replay)
	}
	requireErrorIs(t, err, objectwatch.ErrHistoryUnavailable)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
}

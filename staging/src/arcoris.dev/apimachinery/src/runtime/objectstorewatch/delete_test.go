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

func TestDeletePublishesDeletedChange(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	stream := watchAfter(t, store, testCollection(), created.Revision)

	deleted := deleteObject(t, store, key, created.Revision)
	event := nextEvent(t, stream)

	requireChangedEvent(t, event, watchRequestAfter(t, testCollection(), created.Revision), objectstore.ChangeDeleted, deleted.Revision)
	requireDesiredString(t, event.Change.Before, "created")
	if event.Change.Before.Revision != created.Revision {
		t.Fatalf("before revision = %s; want live revision %s", event.Change.Before.Revision, created.Revision)
	}
	if event.Change.Revision == event.Change.Before.Revision {
		t.Fatalf("delete revision = %s; want tombstone revision distinct from live revision", event.Change.Revision)
	}
	if !event.Change.After.Revision.IsZero() ||
		!event.Change.After.Object.TypeMeta.IsZero() ||
		!event.Change.After.Object.ObjectMeta.IsZero() ||
		!event.Change.After.Object.Desired.IsZero() {
		t.Fatalf("delete after state = %#v; want zero", event.Change.After)
	}
}

func TestDeletePublishesNothingWhenBackendFails(t *testing.T) {
	store := testRuntimeStore(t)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	stream := watchAfter(t, store, testCollection(), created.Revision)

	_, err := store.Delete(context.Background(), key, created.Revision+1)
	if err == nil {
		t.Fatalf("Delete() error = nil; want backend error")
	}

	requireNoEvent(t, stream)
}

func TestDeleteContinuityLossInvalidatesFutureHistoricalWatch(t *testing.T) {
	backend := invalidDeleteResultBackend{Store: testBackend(t)}
	store, err := New(backend)
	requireNoError(t, err)
	key := testKey("system", 1)
	created := createObject(t, store, key, "created")
	stream := watchAfter(t, store, testCollection(), created.Revision)

	deleted, err := store.Delete(context.Background(), key, created.Revision)
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, objectstore.ErrInvalidChange)
	_, err = stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)

	replay, err := store.Watch(context.Background(), watchRequestAfter(t, testCollection(), created.Revision))
	if replay != nil {
		t.Fatalf("stream = %#v; want nil", replay)
	}
	requireErrorIs(t, err, objectwatch.ErrHistoryUnavailable)

	afterGap := watchAfter(t, store, testCollection(), deleted.Revision)
	requireNoEvent(t, afterGap)
}

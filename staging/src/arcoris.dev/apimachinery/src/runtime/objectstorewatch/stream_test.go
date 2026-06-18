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

	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestStreamCloseIsIdempotent(t *testing.T) {
	stream := watchAfter(t, testRuntimeStore(t), testCollection(), 0)

	requireNoError(t, stream.Close())
	requireNoError(t, stream.Close())

	_, err := stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrClosed)
}

func TestStreamCloseUnregistersWatcher(t *testing.T) {
	store := testRuntimeStore(t)
	stream := watchAfter(t, store, testCollection(), 0)
	requireNoError(t, stream.Close())

	createObject(t, store, testKey("system", 1), "created")
	_, err := stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrClosed)
}

func TestStreamDeliveryReturnsDetachedEvents(t *testing.T) {
	store := testRuntimeStore(t)
	created := createObject(t, store, testKey("system", 1), "created")
	stream := watchAfter(t, store, testCollection(), 0)

	event := nextEvent(t, stream)
	event.Change.After.Revision = 99

	replay := watchAfter(t, store, testCollection(), 0)
	replayed := nextEvent(t, replay)
	if replayed.Change.After.Revision != created.Revision {
		t.Fatalf("replayed revision = %s; want %s", replayed.Change.After.Revision, created.Revision)
	}
}

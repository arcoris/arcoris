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

package objectenqueue

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestReflectorSinkIntegrationReplaceWithObjectWorkQueue(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")
	b := testListItem(2, 1, "b")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a, b)))

	requireQueueItems(t, queue, a.Key, b.Key)
}

func TestReflectorSinkIntegrationRelistMissingWithObjectWorkQueue(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")
	b1 := testListItem(2, 1, "b1")
	b2 := testListItem(2, 2, "b2")
	c := testListItem(3, 2, "c")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a, b1)))
	requireQueueItems(t, queue, a.Key, b1.Key)

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 2, b2, c)))

	requireQueueItems(t, queue, b2.Key, c.Key, a.Key)
}

func TestReflectorSinkIntegrationApplyChangeWithObjectWorkQueue(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	sink := newTestReflectorSink(t, queue)
	created := createdChange(t, 1)
	updated := updatedChange(t, 2)
	deleted := deletedChange(t, 3)

	requireNoError(t, sink.ApplyChange(context.Background(), created))
	requireNoError(t, sink.ApplyChange(context.Background(), updated))
	requireNoError(t, sink.ApplyChange(context.Background(), deleted))

	requireQueueItems(t, queue, created.Key, updated.Key, deleted.Key)
}

func requireQueueItems(t testing.TB, queue *objectworkqueue.Queue, want ...objectstore.Key) {
	t.Helper()

	for _, key := range want {
		item, err := queue.Get(context.Background())
		requireNoError(t, err)
		if !item.Key.Equal(key) {
			t.Fatalf("item key = %s; want %v", item.Key, key)
		}
		requireNoError(t, queue.Done(item))
	}
}

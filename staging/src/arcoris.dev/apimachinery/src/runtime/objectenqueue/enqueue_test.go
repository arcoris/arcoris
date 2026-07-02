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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestEnqueuerIntegrationWithObjectWorkQueue(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 1})
	requireNoError(t, err)

	change := createdChange(t, 1)
	handler := mustHandler(t, queue, ChangedObject())

	requireNoError(t, handler.Handle(context.Background(), change))

	item, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, item, objectworkqueue.Item{Key: change.Key})
	requireNoError(t, queue.Done(item))
}

type recordingQueue struct {
	mu       sync.Mutex
	err      error
	contexts []context.Context
	items    []objectworkqueue.Item
}

func (q *recordingQueue) Add(ctx context.Context, item objectworkqueue.Item) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.contexts = append(q.contexts, ctx)
	q.items = append(q.items, item)

	return q.err
}

func (q *recordingQueue) callCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.items)
}

func (q *recordingQueue) requireContexts(t testing.TB, want ...context.Context) {
	t.Helper()

	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.contexts) != len(want) {
		t.Fatalf("contexts = %d; want %d", len(q.contexts), len(want))
	}
	for i := range want {
		if q.contexts[i] != want[i] {
			t.Fatalf("context[%d] = %p; want %p", i, q.contexts[i], want[i])
		}
	}
}

func (q *recordingQueue) requireItems(t testing.TB, want ...objectworkqueue.Item) {
	t.Helper()

	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) != len(want) {
		t.Fatalf("items = %d; want %d", len(q.items), len(want))
	}
	for i := range want {
		requireItem(t, q.items[i], want[i])
	}
}

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

package objectworkqueue

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
}

func newTestQueue(t testing.TB, capacity int) *Queue {
	t.Helper()

	queue, err := New(Options{Capacity: capacity})
	requireNoError(t, err)

	return queue
}

func testItem(id int) Item {
	return Item{Key: testKey(id)}
}

func testKey(id int) objectstore.Key {
	return objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    apiidentity.Group("control.arcoris.dev"),
			Version:  apiidentity.Version("v1"),
			Resource: apiidentity.Resource("workers"),
		},
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace("default"),
			Name:      metaidentity.Name(fmt.Sprintf("worker-%d", id)),
		},
	)
}

func requireItem(t testing.TB, got Item, want Item) {
	t.Helper()

	if !got.Key.Equal(want.Key) {
		t.Fatalf("item = %s; want %s", got.Key, want.Key)
	}
}

func requireStats(t testing.TB, queue *Queue, queued int, processing int) {
	t.Helper()

	stats := queue.Stats()
	if stats.Queued != queued || stats.Processing != processing {
		t.Fatalf("stats = %#v; want queued=%d processing=%d", stats, queued, processing)
	}
	if stats.Queued+stats.Processing > stats.Capacity {
		t.Fatalf("stats = %#v; queued+processing exceeds capacity", stats)
	}
}

func requireInvariants(t testing.TB, queue *Queue) {
	t.Helper()

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.queued > len(queue.items) {
		t.Fatalf("queued=%d tracked=%d; queued exceeds tracked", queue.queued, len(queue.items))
	}
	if (queue.head == nil) != (queue.tail == nil) {
		t.Fatalf("head=%p tail=%p; want both nil or both non-nil", queue.head, queue.tail)
	}

	linked := 0
	seen := make(map[*entry]struct{})
	var prev *entry
	for entry := queue.head; entry != nil; entry = entry.next {
		if _, ok := seen[entry]; ok {
			t.Fatalf("entry %s appears in FIFO more than once", entry.item.Key)
		}
		seen[entry] = struct{}{}
		if entry.prev != prev {
			t.Fatalf("entry %s prev=%p; want %p", entry.item.Key, entry.prev, prev)
		}
		if entry.state != stateQueued {
			t.Fatalf("linked entry %s state=%d; want queued", entry.item.Key, entry.state)
		}
		if entry.dirty {
			t.Fatalf("linked entry %s is dirty", entry.item.Key)
		}
		prev = entry
		linked++
	}
	if prev != queue.tail {
		t.Fatalf("last linked entry=%p tail=%p; want equal", prev, queue.tail)
	}
	if linked != queue.queued {
		t.Fatalf("linked=%d queued=%d; want equal", linked, queue.queued)
	}

	queued := 0
	processing := 0
	for id, entry := range queue.items {
		if keyForItem(entry.item) != id {
			t.Fatalf("entry key = %s; want %s", keyForItem(entry.item), id)
		}

		switch entry.state {
		case stateQueued:
			queued++
			if _, ok := seen[entry]; !ok {
				t.Fatalf("queued entry %s is not linked", id)
			}
			if entry.dirty {
				t.Fatalf("queued entry %s is dirty", id)
			}
		case stateProcessing:
			processing++
			if _, ok := seen[entry]; ok {
				t.Fatalf("processing entry %s is linked", id)
			}
			if entry.prev != nil || entry.next != nil {
				t.Fatalf("processing entry %s has links", id)
			}
		default:
			t.Fatalf("entry %s state = %d; want known state", id, entry.state)
		}
	}

	if queued != queue.queued {
		t.Fatalf("queued=%d queue.queued=%d; want equal", queued, queue.queued)
	}
	if queued+processing != len(queue.items) {
		t.Fatalf("queued=%d processing=%d tracked=%d; want equal", queued, processing, len(queue.items))
	}
	if len(queue.items) > queue.capacity {
		t.Fatalf("tracked=%d capacity=%d; capacity exceeded", len(queue.items), queue.capacity)
	}
}

func waitResult(t testing.TB, ch <-chan error) error {
	t.Helper()

	select {
	case err := <-ch:
		return err
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for result")
		return nil
	}
}

func waitItem(t testing.TB, ch <-chan itemResult) itemResult {
	t.Helper()

	select {
	case result := <-ch:
		return result
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for item")
		return itemResult{}
	}
}

func requireClosed(t testing.TB, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	default:
		t.Fatalf("channel is open; want closed")
	}
}

func withTimeout(t testing.TB) (context.Context, context.CancelFunc) {
	t.Helper()

	return context.WithTimeout(context.Background(), 2*time.Second)
}

type itemResult struct {
	item Item
	err  error
}

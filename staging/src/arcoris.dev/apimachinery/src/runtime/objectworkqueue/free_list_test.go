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
	"testing"
)

func TestFreeListReusesEntriesAfterCleanDone(t *testing.T) {
	queue := newTestQueue(t, 1)
	first := testItem(1)
	second := testItem(2)

	requireNoError(t, queue.Add(context.Background(), first))
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, first)
	requireNoError(t, queue.Done(got))

	released := requireOnlyFreeEntry(t, queue)
	requireStats(t, queue, 0, 0)

	requireNoError(t, queue.Add(context.Background(), second))

	queue.mu.Lock()
	reused := queue.items[keyForItem(second)]
	freeCount := queue.freeCount
	queue.mu.Unlock()

	if reused != released {
		t.Fatalf("reused entry=%p; want released entry=%p", reused, released)
	}
	if freeCount != 0 {
		t.Fatalf("freeCount=%d; want 0 after reuse", freeCount)
	}
	requireStats(t, queue, 1, 0)
	requireInvariants(t, queue)
}

func TestFreeListReusesDirtyShutdownRemovedEntry(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)

	requireNoError(t, queue.Add(context.Background(), item))
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, item)
	requireNoError(t, queue.Add(context.Background(), item))
	queue.ShutDown()
	requireNoError(t, queue.Done(item))

	released := requireOnlyFreeEntry(t, queue)

	queue.mu.Lock()
	reused := queue.acquireEntryLocked(testItem(2))
	queue.releaseEntryLocked(reused)
	queue.mu.Unlock()

	if reused != released {
		t.Fatalf("reused entry=%p; want dirty shutdown entry=%p", reused, released)
	}
	requireStats(t, queue, 0, 0)
	requireInvariants(t, queue)
}

func TestFreeListDoesNotExposeFreeEntriesInStats(t *testing.T) {
	queue := newTestQueue(t, 2)
	item := testItem(1)

	requireNoError(t, queue.Add(context.Background(), item))
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireNoError(t, queue.Done(got))

	stats := queue.Stats()
	if stats.Queued != 0 || stats.Processing != 0 || stats.Capacity != 2 || stats.ShutDown {
		t.Fatalf("stats=%#v; want empty live queue with capacity 2", stats)
	}
	if queue.Len() != 0 {
		t.Fatalf("Len()=%d; want 0", queue.Len())
	}
	requireOnlyFreeEntry(t, queue)
	requireInvariants(t, queue)
}

func TestFreeListBoundedByCapacityAfterHighChurn(t *testing.T) {
	queue := newTestQueue(t, 3)

	for round := 0; round < 250; round++ {
		items := []Item{
			testItem(round*3 + 1),
			testItem(round*3 + 2),
			testItem(round*3 + 3),
		}
		for _, item := range items {
			requireNoError(t, queue.Add(context.Background(), item))
		}
		for _, want := range items {
			got, err := queue.Get(context.Background())
			requireNoError(t, err)
			requireItem(t, got, want)
			requireNoError(t, queue.Done(got))
		}
		requireInvariants(t, queue)
	}

	queue.mu.Lock()
	freeCount := queue.freeCount
	tracked := len(queue.items)
	capacity := queue.capacity
	queue.mu.Unlock()

	if freeCount != capacity {
		t.Fatalf("freeCount=%d; want capacity high-watermark %d", freeCount, capacity)
	}
	if tracked != 0 {
		t.Fatalf("tracked=%d; want 0", tracked)
	}
	requireStats(t, queue, 0, 0)
	requireInvariants(t, queue)
}

func requireOnlyFreeEntry(t testing.TB, queue *Queue) *entry {
	t.Helper()

	queue.mu.Lock()
	defer queue.mu.Unlock()

	if queue.free == nil {
		t.Fatalf("free list is empty")
	}
	if queue.free.next != nil {
		t.Fatalf("free list has more than one entry")
	}
	if queue.freeCount != 1 {
		t.Fatalf("freeCount=%d; want 1", queue.freeCount)
	}
	return queue.free
}

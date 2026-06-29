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

func TestNewRejectsInvalidCapacity(t *testing.T) {
	for _, capacity := range []int{0, -1} {
		_, err := New(Options{Capacity: capacity})
		requireErrorIs(t, err, ErrInvalidCapacity)
	}
}

func TestNewAcceptsPositiveCapacity(t *testing.T) {
	queue, err := New(Options{Capacity: 2})
	requireNoError(t, err)

	stats := queue.Stats()
	if stats.Capacity != 2 {
		t.Fatalf("capacity = %d; want 2", stats.Capacity)
	}
}

func TestQueueStartsEmpty(t *testing.T) {
	queue := newTestQueue(t, 3)

	stats := queue.Stats()
	if stats.Queued != 0 || stats.Processing != 0 {
		t.Fatalf("stats = %#v; want empty", stats)
	}
}

func TestNilQueueBehavior(t *testing.T) {
	var queue *Queue

	requireErrorIs(t, queue.Add(context.Background(), testItem(1)), ErrInvalidQueue)
	requireErrorIs(t, queue.TryAdd(testItem(1)), ErrInvalidQueue)
	_, err := queue.Get(context.Background())
	requireErrorIs(t, err, ErrInvalidQueue)
	requireErrorIs(t, queue.Done(testItem(1)), ErrInvalidQueue)
	if queue.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", queue.Len())
	}
	if stats := queue.Stats(); stats != (Stats{}) {
		t.Fatalf("Stats() = %#v; want zero", stats)
	}
	if !queue.IsShutDown() {
		t.Fatalf("IsShutDown() = false; want true")
	}
	queue.ShutDown()
}

func TestLenReportsQueuedCountOnly(t *testing.T) {
	queue := newTestQueue(t, 2)
	requireNoError(t, queue.Add(context.Background(), testItem(1)))
	requireNoError(t, queue.Add(context.Background(), testItem(2)))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	if queue.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", queue.Len())
	}
}

func TestQueuePrivateCountersRequireLock(t *testing.T) {
	queue := newTestQueue(t, 3)
	requireNoError(t, queue.TryAdd(testItem(1)))

	queue.mu.Lock()
	tracked := queue.trackedLocked()
	queued := queue.queuedLocked()
	processing := queue.processingLocked()
	queue.mu.Unlock()

	if tracked != 1 || queued != 1 || processing != 0 {
		t.Fatalf("tracked=%d queued=%d processing=%d; want 1,1,0", tracked, queued, processing)
	}
}

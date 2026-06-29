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

func TestTryAddEnqueuesNewItem(t *testing.T) {
	queue := newTestQueue(t, 2)

	requireNoError(t, queue.TryAdd(testItem(1)))

	if queue.Len() != 1 {
		t.Fatalf("Len() = %d; want 1", queue.Len())
	}
	requireStats(t, queue, 1, 0)
}

func TestTryAddDuplicateQueuedItemDoesNotIncreaseLen(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)

	requireNoError(t, queue.TryAdd(item))
	requireNoError(t, queue.TryAdd(item))

	requireStats(t, queue, 1, 0)
}

func TestTryAddDuplicateProcessingItemMarksDirty(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.TryAdd(item))
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, item)

	requireNoError(t, queue.TryAdd(item))
	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 1, 0)
	got, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, item)
}

func TestTryAddDuplicateProcessingDirtyItemDoesNotEnqueueMultipleCopies(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.TryAdd(item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	requireNoError(t, queue.TryAdd(item))
	requireNoError(t, queue.TryAdd(item))
	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 1, 0)
	_, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireStats(t, queue, 0, 1)
}

func TestTryAddNewItemWhenFullReturnsFull(t *testing.T) {
	queue := newTestQueue(t, 1)
	requireNoError(t, queue.TryAdd(testItem(1)))

	requireErrorIs(t, queue.TryAdd(testItem(2)), ErrFull)
}

func TestTryAddDuplicateTrackedItemWhenFullSucceeds(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.TryAdd(item))

	requireNoError(t, queue.TryAdd(item))
}

func TestAddEnqueuesNewItemAndDuplicateQueuedItemSucceeds(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)

	requireNoError(t, queue.Add(context.Background(), item))
	requireNoError(t, queue.Add(context.Background(), item))

	requireStats(t, queue, 1, 0)
}

func TestAddDuplicateProcessingItemMarksDirty(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	requireNoError(t, queue.Add(context.Background(), item))
	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 1, 0)
}

func TestAddWaitsWhenFullAndSucceedsAfterDoneFreesCapacity(t *testing.T) {
	queue := newTestQueue(t, 1)
	first := testItem(1)
	second := testItem(2)
	requireNoError(t, queue.Add(context.Background(), first))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	result := make(chan error, 1)
	go func() {
		result <- queue.Add(context.Background(), second)
	}()

	requireNoError(t, queue.Done(first))
	requireNoError(t, waitResult(t, result))
	requireStats(t, queue, 1, 0)
}

func TestAddReturnsContextErrorWhenCancelledWhileWaiting(t *testing.T) {
	queue := newTestQueue(t, 1)
	first := testItem(1)
	second := testItem(2)
	requireNoError(t, queue.Add(context.Background(), first))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())

	result := make(chan error, 1)
	go func() {
		result <- queue.Add(ctx, second)
	}()
	cancel()

	requireErrorIs(t, waitResult(t, result), context.Canceled)
}

func TestAddReturnsShutDownWhenQueueShutsDownWhileWaiting(t *testing.T) {
	queue := newTestQueue(t, 1)
	first := testItem(1)
	second := testItem(2)
	requireNoError(t, queue.Add(context.Background(), first))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	result := make(chan error, 1)
	go func() {
		result <- queue.Add(context.Background(), second)
	}()
	queue.ShutDown()

	requireErrorIs(t, waitResult(t, result), ErrShutDown)
}

func TestAddAndTryAddAfterShutDownReturnShutDown(t *testing.T) {
	queue := newTestQueue(t, 1)
	queue.ShutDown()

	requireErrorIs(t, queue.Add(context.Background(), testItem(1)), ErrShutDown)
	requireErrorIs(t, queue.TryAdd(testItem(1)), ErrShutDown)
}

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

func TestDoneForCleanProcessingItemRemovesTrackingAndFreesCapacity(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 0, 0)
	requireNoError(t, queue.TryAdd(testItem(2)))
}

func TestDoneForUnknownQueuedAndDuplicateItemsReturnsNotProcessing(t *testing.T) {
	queue := newTestQueue(t, 2)
	queued := testItem(1)
	unknown := testItem(2)
	requireNoError(t, queue.Add(context.Background(), queued))

	requireErrorIs(t, queue.Done(unknown), ErrNotProcessing)
	requireErrorIs(t, queue.Done(queued), ErrNotProcessing)
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireNoError(t, queue.Done(queued))
	requireErrorIs(t, queue.Done(queued), ErrNotProcessing)
}

func TestDoneForDirtyProcessingItemRequeuesIt(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireNoError(t, queue.Add(context.Background(), item))

	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 1, 0)
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, item)
}

func TestRepeatedAddWhileProcessingRequeuesOnlyOneCopy(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireNoError(t, queue.Add(context.Background(), item))
	requireNoError(t, queue.Add(context.Background(), item))

	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 1, 0)
	_, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireStats(t, queue, 0, 1)
}

func TestDirtyRequeueAppendsToFIFOBack(t *testing.T) {
	queue := newTestQueue(t, 3)
	first := testItem(1)
	second := testItem(2)
	requireNoError(t, queue.Add(context.Background(), first))
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, first)
	requireNoError(t, queue.Add(context.Background(), second))
	requireNoError(t, queue.Add(context.Background(), first))
	requireNoError(t, queue.Done(first))

	got, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, second)
	got, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, first)
}

func TestDoneDirtyItemAfterShutDownRemovesInsteadOfRequeueing(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireNoError(t, queue.Add(context.Background(), item))
	queue.ShutDown()

	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 0, 0)
	_, err = queue.Get(context.Background())
	requireErrorIs(t, err, ErrShutDown)
}

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

func TestGetReturnsQueuedItemAndMarksProcessing(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))

	got, err := queue.Get(context.Background())
	requireNoError(t, err)

	requireItem(t, got, item)
	requireStats(t, queue, 0, 1)
}

func TestGetReturnsDistinctItemsInFIFOOrder(t *testing.T) {
	queue := newTestQueue(t, 3)
	first := testItem(1)
	second := testItem(2)
	third := testItem(3)
	requireNoError(t, queue.Add(context.Background(), first))
	requireNoError(t, queue.Add(context.Background(), second))
	requireNoError(t, queue.Add(context.Background(), third))

	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, first)
	got, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, second)
	got, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, third)
}

func TestGetBlocksUntilAddEnqueuesItem(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	result := make(chan itemResult, 1)
	ctx, cancel := withTimeout(t)
	defer cancel()

	go func() {
		got, err := queue.Get(ctx)
		result <- itemResult{item: got, err: err}
	}()
	requireNoError(t, queue.Add(context.Background(), item))

	got := waitItem(t, result)
	requireNoError(t, got.err)
	requireItem(t, got.item, item)
}

func TestGetReturnsContextErrorWhenCancelledWhileWaiting(t *testing.T) {
	queue := newTestQueue(t, 1)
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan itemResult, 1)

	go func() {
		got, err := queue.Get(ctx)
		result <- itemResult{item: got, err: err}
	}()
	cancel()

	requireErrorIs(t, waitItem(t, result).err, context.Canceled)
}

func TestGetAfterShutDownAndEmptyQueueReturnsShutDown(t *testing.T) {
	queue := newTestQueue(t, 1)
	queue.ShutDown()

	_, err := queue.Get(context.Background())
	requireErrorIs(t, err, ErrShutDown)
}

func TestGetAfterShutDownDrainsQueuedItemsBeforeShutDown(t *testing.T) {
	queue := newTestQueue(t, 2)
	first := testItem(1)
	second := testItem(2)
	requireNoError(t, queue.Add(context.Background(), first))
	requireNoError(t, queue.Add(context.Background(), second))
	queue.ShutDown()

	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, first)
	got, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, second)
	_, err = queue.Get(context.Background())
	requireErrorIs(t, err, ErrShutDown)
}

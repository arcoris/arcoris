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

func TestSignalNotEmptyLockedSkipsChannelWithoutWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notEmpty
	queue.signalNotEmptyLocked()
	replacement := queue.notEmpty
	queue.mu.Unlock()

	if replacement == ch {
		return
	}
	t.Fatalf("notEmpty channel was replaced without waiters")
}

func TestSignalNotEmptyLockedClosesAndReplacesChannelWithWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notEmpty
	queue.notEmptyWaiters = 1
	queue.signalNotEmptyLocked()
	replacement := queue.notEmpty
	queue.mu.Unlock()

	requireClosed(t, ch)
	if replacement == ch {
		t.Fatalf("notEmpty channel was not replaced")
	}
}

func TestSignalNotFullLockedSkipsChannelWithoutWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notFull
	queue.signalNotFullLocked()
	replacement := queue.notFull
	queue.mu.Unlock()

	if replacement == ch {
		return
	}
	t.Fatalf("notFull channel was replaced without waiters")
}

func TestSignalNotFullLockedClosesAndReplacesChannelWithWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	ch := queue.notFull
	queue.notFullWaiters = 1
	queue.signalNotFullLocked()
	replacement := queue.notFull
	queue.mu.Unlock()

	requireClosed(t, ch)
	if replacement == ch {
		t.Fatalf("notFull channel was not replaced")
	}
}

func TestSignalAllLockedClosesAndReplacesBothChannelsWithWaiters(t *testing.T) {
	queue := newTestQueue(t, 1)

	queue.mu.Lock()
	notEmpty := queue.notEmpty
	notFull := queue.notFull
	queue.notEmptyWaiters = 1
	queue.notFullWaiters = 1
	queue.signalAllLocked()
	newNotEmpty := queue.notEmpty
	newNotFull := queue.notFull
	queue.mu.Unlock()

	requireClosed(t, notEmpty)
	requireClosed(t, notFull)
	if newNotEmpty == notEmpty || newNotFull == notFull {
		t.Fatalf("signalAllLocked did not replace both channels")
	}
}

func TestNotEmptyWaiterCountReturnsToZeroAfterSignal(t *testing.T) {
	queue := newTestQueue(t, 1)
	result := make(chan itemResult, 1)
	go func() {
		item, err := queue.Get(context.Background())
		result <- itemResult{item: item, err: err}
	}()
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 1
	})

	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	got := waitItem(t, result)
	requireNoError(t, got.err)
	requireItem(t, got.item, item)
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 0
	})
}

func TestNotFullWaiterCountReturnsToZeroAfterSignal(t *testing.T) {
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
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notFullWaiters == 1
	})

	requireNoError(t, queue.Done(first))
	requireNoError(t, waitResult(t, result))
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notFullWaiters == 0
	})
}

func TestNotEmptyWaiterCountReturnsToZeroAfterContextCancel(t *testing.T) {
	queue := newTestQueue(t, 1)
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan itemResult, 1)
	go func() {
		item, err := queue.Get(ctx)
		result <- itemResult{item: item, err: err}
	}()
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 1
	})

	cancel()
	requireErrorIs(t, waitItem(t, result).err, context.Canceled)
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 0
	})
}

func TestNotFullWaiterCountReturnsToZeroAfterContextCancel(t *testing.T) {
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
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notFullWaiters == 1
	})

	cancel()
	requireErrorIs(t, waitResult(t, result), context.Canceled)
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notFullWaiters == 0
	})
}

func TestManyBlockedGetWaitersWakeOnShutDown(t *testing.T) {
	queue := newTestQueue(t, 1)
	results := make(chan error, 8)
	for range 8 {
		go func() {
			_, err := queue.Get(context.Background())
			results <- err
		}()
	}
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 8
	})

	queue.ShutDown()
	for range 8 {
		requireErrorIs(t, waitResult(t, results), ErrShutDown)
	}
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 0
	})
}

func TestManyBlockedAddWaitersWakeOnShutDown(t *testing.T) {
	queue := newTestQueue(t, 1)
	first := testItem(1)
	requireNoError(t, queue.Add(context.Background(), first))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)
	results := make(chan error, 8)
	for i := range 8 {
		item := testItem(i + 2)
		go func() {
			results <- queue.Add(context.Background(), item)
		}()
	}
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notFullWaiters == 8
	})

	queue.ShutDown()
	for range 8 {
		requireErrorIs(t, waitResult(t, results), ErrShutDown)
	}
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notFullWaiters == 0
	})
}

func TestRepeatedSignalsWithSlowWaiterDoNotBreakFutureWaiters(t *testing.T) {
	queue := newTestQueue(t, 2)

	queue.mu.Lock()
	ch := queue.notEmpty
	queue.notEmptyWaiters = 1
	queue.signalNotEmptyLocked()
	queue.signalNotEmptyLocked()
	queue.notEmptyWaiters--
	queue.mu.Unlock()
	requireClosed(t, ch)

	result := make(chan itemResult, 1)
	go func() {
		item, err := queue.Get(context.Background())
		result <- itemResult{item: item, err: err}
	}()
	waitUntil(t, func() bool {
		queue.mu.Lock()
		defer queue.mu.Unlock()
		return queue.notEmptyWaiters == 1
	})

	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	got := waitItem(t, result)
	requireNoError(t, got.err)
	requireItem(t, got.item, item)
}

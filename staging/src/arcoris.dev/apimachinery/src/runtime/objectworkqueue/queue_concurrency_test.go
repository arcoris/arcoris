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
	"sync"
	"testing"
)

func TestConcurrentTryAddSameItemQueuesOneItem(t *testing.T) {
	queue := newTestQueue(t, 64)
	item := testItem(1)
	errs := runParallel(64, func() error {
		return queue.TryAdd(item)
	})

	requireNoParallelErrors(t, errs)
	requireStats(t, queue, 1, 0)
}

func TestConcurrentAddSameProcessingItemRequeuesOnce(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	errs := runParallel(64, func() error {
		return queue.Add(context.Background(), item)
	})
	requireNoParallelErrors(t, errs)
	requireNoError(t, queue.Done(item))

	requireStats(t, queue, 1, 0)
	_, err = queue.Get(context.Background())
	requireNoError(t, err)
	requireStats(t, queue, 0, 1)
}

func TestConcurrentTryAddDistinctItemsUpToCapacity(t *testing.T) {
	const total = 32
	queue := newTestQueue(t, total)
	errs := runParallelIndexed(total, func(i int) error {
		return queue.TryAdd(testItem(i))
	})

	requireNoParallelErrors(t, errs)
	requireStats(t, queue, total, 0)
}

func TestConcurrentWorkersGetAndDoneItems(t *testing.T) {
	const total = 32
	queue := newTestQueue(t, total)
	for i := range total {
		requireNoError(t, queue.Add(context.Background(), testItem(i)))
	}

	errs := runParallel(total, func() error {
		item, err := queue.Get(context.Background())
		if err != nil {
			return err
		}
		return queue.Done(item)
	})

	requireNoParallelErrors(t, errs)
	requireStats(t, queue, 0, 0)
}

func TestBlockedAddWaitersWakeWhenCapacityIsFreed(t *testing.T) {
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
}

func TestBlockedGetWaitersWakeWhenItemsAreAdded(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	result := make(chan itemResult, 1)

	go func() {
		got, err := queue.Get(context.Background())
		result <- itemResult{item: got, err: err}
	}()
	requireNoError(t, queue.Add(context.Background(), item))

	got := waitItem(t, result)
	requireNoError(t, got.err)
	requireItem(t, got.item, item)
}

func TestBlockedAddAndGetWaitersWakeOnShutDown(t *testing.T) {
	queue := newTestQueue(t, 1)
	first := testItem(1)
	second := testItem(2)
	requireNoError(t, queue.Add(context.Background(), first))
	_, err := queue.Get(context.Background())
	requireNoError(t, err)

	addResult := make(chan error, 1)
	getResult := make(chan itemResult, 1)
	go func() {
		addResult <- queue.Add(context.Background(), second)
	}()
	emptyQueue := newTestQueue(t, 1)
	go func() {
		got, err := emptyQueue.Get(context.Background())
		getResult <- itemResult{item: got, err: err}
	}()

	queue.ShutDown()
	emptyQueue.ShutDown()

	requireErrorIs(t, waitResult(t, addResult), ErrShutDown)
	requireErrorIs(t, waitItem(t, getResult).err, ErrShutDown)
}

func TestConcurrentAddAndDoneLeavesItemQueued(t *testing.T) {
	for i := range 1000 {
		queue := newTestQueue(t, 2)
		item := testItem(i)
		requireNoError(t, queue.Add(context.Background(), item))
		_, err := queue.Get(context.Background())
		requireNoError(t, err)

		start := make(chan struct{})
		errs := make(chan error, 2)
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			<-start
			errs <- queue.Add(context.Background(), item)
		}()
		go func() {
			defer wg.Done()
			<-start
			errs <- queue.Done(item)
		}()
		close(start)
		wg.Wait()
		close(errs)

		for err := range errs {
			requireNoError(t, err)
		}
		requireStats(t, queue, 1, 0)
		requireInvariants(t, queue)
		got, err := queue.Get(context.Background())
		requireNoError(t, err)
		requireItem(t, got, item)
	}
}

func runParallel(count int, fn func() error) []error {
	errs := make([]error, count)
	var wg sync.WaitGroup
	wg.Add(count)
	for i := range count {
		go func() {
			defer wg.Done()
			errs[i] = fn()
		}()
	}
	wg.Wait()
	return errs
}

func runParallelIndexed(count int, fn func(int) error) []error {
	errs := make([]error, count)
	var wg sync.WaitGroup
	wg.Add(count)
	for i := range count {
		go func() {
			defer wg.Done()
			errs[i] = fn(i)
		}()
	}
	wg.Wait()
	return errs
}

func requireNoParallelErrors(t testing.TB, errs []error) {
	t.Helper()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("error[%d] = %v", i, err)
		}
	}
}

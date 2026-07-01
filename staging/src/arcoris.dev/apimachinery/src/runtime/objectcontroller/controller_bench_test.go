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

package objectcontroller

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
	"arcoris.dev/snapshot"
)

const benchmarkRunBatchSize = 1024

func BenchmarkProcessItemNoopReconciler(b *testing.B) {
	controller := newBenchProcessController(b, objectreconciler.Success(), nil)
	item := testItem(1)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := controller.processItem(ctx, item); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkProcessItemSourceError(b *testing.B) {
	sourceErr := errors.New("source failed")
	controller := newBenchProcessController(b, objectreconciler.Success(), sourceErr)
	item := testItem(1)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := controller.processItem(ctx, item); err != sourceErr {
			b.Fatalf("error = %v; want %v", err, sourceErr)
		}
	}
}

func BenchmarkProcessItemReconcileError(b *testing.B) {
	reconcileErr := errors.New("reconcile failed")
	controller := newBenchProcessController(b, objectreconciler.Failure(reconcileErr), nil)
	item := testItem(1)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := controller.processItem(ctx, item); err != reconcileErr {
			b.Fatalf("error = %v; want %v", err, reconcileErr)
		}
	}
}

func BenchmarkRunWithRealObjectWorkQueueSingleWorker(b *testing.B) {
	benchmarkRunWithRealObjectWorkQueue(b, 1)
}

func BenchmarkRunWithRealObjectWorkQueueParallelWorkers(b *testing.B) {
	for _, workers := range []int{1, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
			benchmarkRunWithRealObjectWorkQueue(b, workers)
		})
	}
}

func BenchmarkRunWithFakeQueueSingleWorker(b *testing.B) {
	benchmarkRunWithFakeQueue(b, 1)
}

func BenchmarkRunWithFakeQueueParallelWorkers(b *testing.B) {
	for _, workers := range []int{1, 2, 4, 8, 16} {
		b.Run(fmt.Sprintf("workers=%d", workers), func(b *testing.B) {
			benchmarkRunWithFakeQueue(b, workers)
		})
	}
}

func BenchmarkRunQueueShutdownOnly(b *testing.B) {
	snap := testSnapshot(b, 1)
	ctx := context.Background()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		queue := &benchQueue{}
		controller := newBenchController(b, Options{Workers: 1}, queue, snap, &benchReconciler{result: objectreconciler.Success()})
		if err := controller.Run(ctx); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRunContextCancelledWhileWaiting(b *testing.B) {
	snap := testSnapshot(b, 1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		controller := newBenchController(b, Options{Workers: 1}, &benchBlockingQueue{}, snap, &benchReconciler{result: objectreconciler.Success()})
		if err := controller.Run(ctx); !errors.Is(err, context.Canceled) {
			b.Fatalf("error = %v; want context.Canceled", err)
		}
	}
}

func benchmarkRunWithRealObjectWorkQueue(b *testing.B, workers int) {
	snap := testSnapshot(b, 1)
	ctx := context.Background()
	processed := 0

	b.ReportAllocs()
	b.ResetTimer()

	for processed < b.N {
		n := min(benchmarkRunBatchSize, b.N-processed)
		b.StopTimer()
		queue := newBenchRealQueue(b, n)
		queue.ShutDown()
		reconciler := &benchReconciler{result: objectreconciler.Success()}
		controller, err := New(Options{Workers: workers}, queue, &benchSource{snapshot: snap}, reconciler)
		if err != nil {
			b.Fatal(err)
		}
		b.StartTimer()
		err = controller.Run(ctx)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		if calls := reconciler.callCount(); calls != int64(n) {
			b.Fatalf("reconciler calls = %d; want %d", calls, n)
		}
		stats := queue.Stats()
		if stats.Queued != 0 || stats.Processing != 0 {
			b.Fatalf("queue stats = %#v; want drained queue", stats)
		}
		processed += n
	}
}

func benchmarkRunWithFakeQueue(b *testing.B, workers int) {
	snap := testSnapshot(b, 1)
	ctx := context.Background()
	processed := 0

	b.ReportAllocs()
	b.ResetTimer()

	for processed < b.N {
		n := min(benchmarkRunBatchSize, b.N-processed)
		b.StopTimer()
		queue := newBenchQueue(n)
		reconciler := &benchReconciler{result: objectreconciler.Success()}
		controller := newBenchController(b, Options{Workers: workers}, queue, snap, reconciler)
		b.StartTimer()
		err := controller.Run(ctx)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
		if calls := reconciler.callCount(); calls != int64(n) {
			b.Fatalf("reconciler calls = %d; want %d", calls, n)
		}
		if done := queue.doneCount(); done != int64(n) {
			b.Fatalf("Done calls = %d; want %d", done, n)
		}
		processed += n
	}
}

func newBenchProcessController(b testing.TB, result objectreconciler.Result, sourceErr error) *Controller {
	b.Helper()

	snap := testSnapshot(b, 1)
	reconciler := &benchReconciler{result: result}
	controller, err := New(
		Options{Workers: 1},
		&benchDoneQueue{},
		&benchSource{snapshot: snap, err: sourceErr},
		reconciler,
	)
	if err != nil {
		b.Fatal(err)
	}
	return controller
}

func newBenchController(
	b testing.TB,
	opts Options,
	queue Queue,
	snap benchSnapshot,
	reconciler objectreconciler.Reconciler,
) *Controller {
	b.Helper()

	controller, err := New(opts, queue, &benchSource{snapshot: snap}, reconciler)
	if err != nil {
		b.Fatal(err)
	}
	return controller
}

func newBenchRealQueue(b testing.TB, n int) *objectworkqueue.Queue {
	b.Helper()

	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: n})
	if err != nil {
		b.Fatal(err)
	}
	for i := 0; i < n; i++ {
		if err := queue.Add(context.Background(), testItem(i+1)); err != nil {
			b.Fatal(err)
		}
	}
	return queue
}

func newBenchQueue(n int) *benchQueue {
	queue := &benchQueue{
		items: make([]objectworkqueue.Item, n),
	}
	for i := range queue.items {
		queue.items[i] = testItem(i + 1)
	}
	return queue
}

type benchSnapshot = snapshot.Snapshot[objectstore.Revision, objectcache.View]

type benchSource struct {
	snapshot benchSnapshot
	err      error
}

func (s *benchSource) ReadSnapshot() (benchSnapshot, error) {
	if s.err != nil {
		return benchSnapshot{}, s.err
	}
	return s.snapshot, nil
}

type benchReconciler struct {
	result objectreconciler.Result
	calls  atomic.Int64
}

func (r *benchReconciler) Reconcile(context.Context, objectreconciler.Snapshot) objectreconciler.Result {
	r.calls.Add(1)
	return r.result
}

func (r *benchReconciler) callCount() int64 {
	return r.calls.Load()
}

type benchDoneQueue struct {
	done atomic.Int64
}

func (q *benchDoneQueue) Get(context.Context) (objectworkqueue.Item, error) {
	return objectworkqueue.Item{}, objectworkqueue.ErrShutDown
}

func (q *benchDoneQueue) Done(objectworkqueue.Item) error {
	q.done.Add(1)
	return nil
}

type benchQueue struct {
	next  atomic.Uint64
	done  atomic.Int64
	items []objectworkqueue.Item
}

func (q *benchQueue) Get(context.Context) (objectworkqueue.Item, error) {
	next := int(q.next.Add(1))
	if next > len(q.items) {
		return objectworkqueue.Item{}, objectworkqueue.ErrShutDown
	}
	return q.items[next-1], nil
}

func (q *benchQueue) Done(objectworkqueue.Item) error {
	q.done.Add(1)
	return nil
}

func (q *benchQueue) doneCount() int64 {
	return q.done.Load()
}

type benchBlockingQueue struct{}

func (*benchBlockingQueue) Get(ctx context.Context) (objectworkqueue.Item, error) {
	<-ctx.Done()
	return objectworkqueue.Item{}, ctx.Err()
}

func (*benchBlockingQueue) Done(objectworkqueue.Item) error {
	return nil
}

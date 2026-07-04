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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestRunNilControllerReturnsErrInvalidController(t *testing.T) {
	var controller *Controller

	requireErrorIs(t, controller.Run(context.Background()), ErrInvalidController)
}

func TestRunPanicsOnNilContext(t *testing.T) {
	controller, err := New(Options{Workers: 1}, &recordingQueue{}, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)

	defer func() {
		if got := recover(); got != "objectcontroller: nil context" {
			t.Fatalf("recover() = %#v; want objectcontroller: nil context", got)
		}
	}()

	_ = controller.Run(nil)
}

func TestRunRejectsConcurrentRun(t *testing.T) {
	queue := &blockingQueue{started: make(chan struct{}, 1)}
	controller, err := New(Options{Workers: 1}, queue, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := make(chan error, 1)
	go func() {
		result <- controller.Run(ctx)
	}()
	<-queue.started

	requireErrorIs(t, controller.Run(context.Background()), ErrAlreadyRunning)
	cancel()
	requireErrorIs(t, waitError(t, result), context.Canceled)
}

func TestRunAllowsSequentialRunsAfterCompletion(t *testing.T) {
	controller, err := New(Options{Workers: 1}, &recordingQueue{}, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)

	requireNoError(t, controller.Run(context.Background()))
	requireNoError(t, controller.Run(context.Background()))
}

func TestRunStartsConfiguredWorkers(t *testing.T) {
	queue := &blockingQueue{started: make(chan struct{}, 3)}
	controller, err := New(Options{Workers: 3}, queue, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := make(chan error, 1)
	go func() {
		result <- controller.Run(ctx)
	}()
	waitUntil(t, func() bool {
		return queue.getCount() == 3
	})

	cancel()
	requireErrorIs(t, waitError(t, result), context.Canceled)
}

func TestRunExitsCleanlyWhenQueueShutsDown(t *testing.T) {
	controller, err := New(Options{Workers: 2}, &recordingQueue{}, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)

	requireNoError(t, controller.Run(context.Background()))
}

func TestRunReturnsContextErrorWhenCancelledWhileWaiting(t *testing.T) {
	queue := &blockingQueue{started: make(chan struct{}, 1)}
	controller, err := New(Options{Workers: 1}, queue, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())

	result := make(chan error, 1)
	go func() {
		result <- controller.Run(ctx)
	}()
	<-queue.started
	cancel()

	requireErrorIs(t, waitError(t, result), context.Canceled)
}

func TestRunCancelsOtherWorkersAfterFatalError(t *testing.T) {
	fatalErr := errors.New("fatal get")
	queue := &fatalThenBlockingQueue{
		fatal:         fatalErr,
		secondStarted: make(chan struct{}, 1),
	}
	controller, err := New(Options{Workers: 2}, queue, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)

	err = controller.Run(context.Background())

	requireErrorSame(t, err, fatalErr)
	if queue.cancelledCount() == 0 {
		t.Fatal("other worker did not observe cancellation")
	}
}

func TestRunWaitsForWorkersBeforeReturning(t *testing.T) {
	fatalErr := errors.New("fatal get")
	release := make(chan struct{})
	queue := &fatalThenBlockingQueue{
		fatal:         fatalErr,
		secondStarted: make(chan struct{}, 1),
		release:       release,
	}
	controller, err := New(Options{Workers: 2}, queue, &fakeSnapshotSource{}, &fakeReconciler{})
	requireNoError(t, err)
	result := make(chan error, 1)

	go func() {
		result <- controller.Run(context.Background())
	}()
	<-queue.secondStarted
	waitUntil(t, func() bool {
		return queue.cancelledCount() == 1
	})
	assertNoErrorYet(t, result)
	close(release)

	requireErrorSame(t, waitError(t, result), fatalErr)
}

func TestRunDoesNotShutDownQueue(t *testing.T) {
	queue := &blockingQueue{}
	controller, err := New(Options{Workers: 1}, queue, &fakeSnapshotSource{}, objectreconciler.ReconcileFunc(func(context.Context, objectreconciler.Request, objectreconciler.Snapshot) objectreconciler.Result {
		return objectreconciler.Success()
	}))
	requireNoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	requireErrorIs(t, controller.Run(ctx), context.Canceled)
	if queue.shutDownCount() != 0 {
		t.Fatalf("ShutDown calls = %d; want 0", queue.shutDownCount())
	}
}

func TestRunReturnsReconcileErrorAfterDone(t *testing.T) {
	reconcileErr := errors.New("reconcile failed")
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1)}}
	controller, err := New(
		Options{Workers: 1},
		queue,
		&fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		&fakeReconciler{result: objectreconciler.Failure(reconcileErr)},
	)
	requireNoError(t, err)

	requireErrorSame(t, controller.Run(context.Background()), reconcileErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

func TestRunReturnsDoneErrorAfterSuccessfulReconcile(t *testing.T) {
	doneErr := errors.New("done failed")
	queue := &recordingQueue{
		items:   []objectworkqueue.Item{testItem(1)},
		doneErr: doneErr,
	}
	controller, err := New(
		Options{Workers: 1},
		queue,
		&fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		&fakeReconciler{result: objectreconciler.Success()},
	)
	requireNoError(t, err)

	requireErrorSame(t, controller.Run(context.Background()), doneErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

func TestRunReturnsJoinedReconcileAndDoneErrors(t *testing.T) {
	reconcileErr := errors.New("reconcile failed")
	doneErr := errors.New("done failed")
	queue := &recordingQueue{
		items:   []objectworkqueue.Item{testItem(1)},
		doneErr: doneErr,
	}
	controller, err := New(
		Options{Workers: 1},
		queue,
		&fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		&fakeReconciler{result: objectreconciler.Failure(reconcileErr)},
	)
	requireNoError(t, err)

	err = controller.Run(context.Background())

	requireErrorIs(t, err, reconcileErr)
	requireErrorIs(t, err, doneErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

type blockingQueue struct {
	mu       sync.Mutex
	gets     int
	started  chan struct{}
	done     int
	shutdown int
}

func (q *blockingQueue) Get(ctx context.Context) (objectworkqueue.Item, error) {
	q.mu.Lock()
	q.gets++
	started := q.started
	q.mu.Unlock()

	if started != nil {
		started <- struct{}{}
	}
	<-ctx.Done()
	return objectworkqueue.Item{}, ctx.Err()
}

func (q *blockingQueue) Done(objectworkqueue.Item) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.done++
	return nil
}

func (q *blockingQueue) ShutDown() {
	q.mu.Lock()
	q.shutdown++
	q.mu.Unlock()
}

func (q *blockingQueue) getCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.gets
}

func (q *blockingQueue) shutDownCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.shutdown
}

type fatalThenBlockingQueue struct {
	mu            sync.Mutex
	fatal         error
	release       <-chan struct{}
	secondStarted chan struct{}
	calls         int
	cancelled     int
}

func (q *fatalThenBlockingQueue) Get(ctx context.Context) (objectworkqueue.Item, error) {
	q.mu.Lock()
	q.calls++
	call := q.calls
	fatal := q.fatal
	secondStarted := q.secondStarted
	release := q.release
	q.mu.Unlock()

	if call == 1 {
		return objectworkqueue.Item{}, fatal
	}
	if secondStarted != nil {
		select {
		case secondStarted <- struct{}{}:
		default:
		}
	}
	<-ctx.Done()
	q.mu.Lock()
	q.cancelled++
	q.mu.Unlock()
	if release != nil {
		<-release
	}
	return objectworkqueue.Item{}, ctx.Err()
}

func (q *fatalThenBlockingQueue) Done(objectworkqueue.Item) error {
	return nil
}

func (q *fatalThenBlockingQueue) cancelledCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return q.cancelled
}

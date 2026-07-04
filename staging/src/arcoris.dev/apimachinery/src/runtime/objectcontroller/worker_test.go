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
	"testing"

	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestWorkerReturnsNilOnQueueShutdown(t *testing.T) {
	controller := &Controller{
		queue:      &recordingQueue{getErrs: []error{objectworkqueue.ErrShutDown}},
		source:     &fakeSnapshotSource{},
		reconciler: &fakeReconciler{},
	}

	requireNoError(t, controller.runWorker(context.Background()))
}

func TestWorkerReturnsGetError(t *testing.T) {
	getErr := errors.New("get failed")
	controller := &Controller{
		queue:      &recordingQueue{getErrs: []error{getErr}},
		source:     &fakeSnapshotSource{},
		reconciler: &fakeReconciler{},
	}

	requireErrorSame(t, controller.runWorker(context.Background()), getErr)
}

func TestWorkerProcessesItemAndCallsDone(t *testing.T) {
	item := testItem(1)
	queue := &recordingQueue{items: []objectworkqueue.Item{item}}
	reconciler := &fakeReconciler{result: objectreconciler.Success()}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: reconciler,
	}

	requireNoError(t, controller.runWorker(context.Background()))

	if reconciler.callCount() != 1 {
		t.Fatalf("reconciler calls = %d; want 1", reconciler.callCount())
	}
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
	requests := reconciler.reconcilerRequests()
	if len(requests) != 1 {
		t.Fatalf("reconciler requests = %d; want 1", len(requests))
	}
	if !requests[0].Key.Equal(item.Key) {
		t.Fatalf("request key = %#v; want %#v", requests[0].Key, item.Key)
	}
	doneItems := queue.completedItems()
	if len(doneItems) != 1 {
		t.Fatalf("completed items = %d; want 1", len(doneItems))
	}
	if !doneItems[0].Key.Equal(item.Key) {
		t.Fatalf("completed item key = %#v; want %#v", doneItems[0].Key, item.Key)
	}
}

func TestWorkerContinuesAfterSuccessfulReconcile(t *testing.T) {
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1), testItem(2)}}
	reconciler := &fakeReconciler{result: objectreconciler.Success()}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: reconciler,
	}

	requireNoError(t, controller.runWorker(context.Background()))

	if reconciler.callCount() != 2 {
		t.Fatalf("reconciler calls = %d; want 2", reconciler.callCount())
	}
	if queue.doneCount() != 2 {
		t.Fatalf("Done calls = %d; want 2", queue.doneCount())
	}
}

func TestWorkerReturnsReconcileErrorAfterDone(t *testing.T) {
	reconcileErr := errors.New("reconcile failed")
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1)}}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: &fakeReconciler{result: objectreconciler.Failure(reconcileErr)},
	}

	requireErrorSame(t, controller.runWorker(context.Background()), reconcileErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

func TestWorkerReturnsDoneError(t *testing.T) {
	doneErr := errors.New("done failed")
	queue := &recordingQueue{
		items:   []objectworkqueue.Item{testItem(1)},
		doneErr: doneErr,
	}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: &fakeReconciler{result: objectreconciler.Success()},
	}

	requireErrorSame(t, controller.runWorker(context.Background()), doneErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

func TestWorkerJoinsReconcileAndDoneErrors(t *testing.T) {
	reconcileErr := errors.New("reconcile failed")
	doneErr := errors.New("done failed")
	queue := &recordingQueue{
		items:   []objectworkqueue.Item{testItem(1)},
		doneErr: doneErr,
	}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: &fakeReconciler{result: objectreconciler.Failure(reconcileErr)},
	}

	err := controller.runWorker(context.Background())

	requireErrorIs(t, err, reconcileErr)
	requireErrorIs(t, err, doneErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

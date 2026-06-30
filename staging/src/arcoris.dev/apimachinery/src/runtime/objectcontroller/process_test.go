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

func TestWorkerCallsDoneWhenSnapshotSourceFails(t *testing.T) {
	sourceErr := errors.New("snapshot failed")
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1)}}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{err: sourceErr},
		reconciler: &fakeReconciler{},
	}

	requireErrorSame(t, controller.processItem(context.Background(), testItem(1)), sourceErr)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

func TestWorkerCallsDoneWhenContextCancelledAfterGet(t *testing.T) {
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1)}}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: &fakeReconciler{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	requireErrorIs(t, controller.processItem(ctx, testItem(1)), context.Canceled)
	if queue.doneCount() != 1 {
		t.Fatalf("Done calls = %d; want 1", queue.doneCount())
	}
}

func TestWorkerAttemptsDoneBeforePanicPropagates(t *testing.T) {
	queue := &recordingQueue{items: []objectworkqueue.Item{testItem(1)}}
	controller := &Controller{
		queue:      queue,
		source:     &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: &fakeReconciler{panicValue: "boom"},
	}

	defer func() {
		if got := recover(); got != "boom" {
			t.Fatalf("recover() = %#v; want boom", got)
		}
		if queue.doneCount() != 1 {
			t.Fatalf("Done calls = %d; want 1", queue.doneCount())
		}
	}()

	_ = controller.processItem(context.Background(), testItem(1))
}

func TestProcessItemReturnsDoneErrorWhenReconcileSucceeds(t *testing.T) {
	doneErr := errors.New("done failed")
	queue := &recordingQueue{doneErr: doneErr}
	controller := &Controller{
		queue:  queue,
		source: &fakeSnapshotSource{snapshot: testSnapshot(t, 1)},
		reconciler: objectreconciler.ReconcileFunc(func(context.Context, objectreconciler.Snapshot) objectreconciler.Result {
			return objectreconciler.Success()
		}),
	}

	requireErrorSame(t, controller.processItem(context.Background(), testItem(1)), doneErr)
}

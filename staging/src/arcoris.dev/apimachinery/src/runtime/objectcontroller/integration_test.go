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
	"testing"

	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

var _ Queue = (*objectworkqueue.Queue)(nil)

func TestRunWithRealObjectWorkQueue(t *testing.T) {
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	requireNoError(t, queue.Add(context.Background(), testItem(1)))
	requireNoError(t, queue.Add(context.Background(), testItem(2)))
	queue.ShutDown()

	reconciler := &fakeReconciler{result: objectreconciler.Success()}
	controller, err := New(
		Options{Workers: 2},
		queue,
		&fakeSnapshotSource{snapshot: testSnapshot(t, 7)},
		reconciler,
	)
	requireNoError(t, err)

	requireNoError(t, controller.Run(context.Background()))

	if reconciler.callCount() != 2 {
		t.Fatalf("reconciler calls = %d; want 2", reconciler.callCount())
	}
	stats := queue.Stats()
	if stats.Queued != 0 || stats.Processing != 0 {
		t.Fatalf("queue stats = %#v; want drained queue", stats)
	}
}

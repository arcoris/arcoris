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

package objectreconciler

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/runtime/objectcache"
)

func TestReconcileOncePanicsOnNilContext(t *testing.T) {
	defer func() {
		if got := recover(); got != "objectreconciler: nil context" {
			t.Fatalf("recover() = %#v; want objectreconciler: nil context", got)
		}
	}()

	_ = ReconcileOnce(nil, testRequest(1), &fakeSource{}, &recordingReconciler{})
}

func TestReconcileOnceReturnsCancelledContextBeforeSourceRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	source := &fakeSource{snapshot: readySourceSnapshot(t, 1)}
	reconciler := &recordingReconciler{}

	result := ReconcileOnce(ctx, Request{}, source, reconciler)

	requireErrorIs(t, result.Err, context.Canceled)
	if source.reads != 0 {
		t.Fatalf("source reads = %d; want 0", source.reads)
	}
	if reconciler.calls != 0 {
		t.Fatalf("reconciler calls = %d; want 0", reconciler.calls)
	}
}

func TestReconcileOnceRejectsInvalidRequestBeforeSourceRead(t *testing.T) {
	source := &fakeSource{snapshot: readySourceSnapshot(t, 1)}
	reconciler := &recordingReconciler{}

	result := ReconcileOnce(context.Background(), Request{}, source, reconciler)

	requireErrorIs(t, result.Err, ErrInvalidRequest)
	if source.reads != 0 {
		t.Fatalf("source reads = %d; want 0", source.reads)
	}
	if reconciler.calls != 0 {
		t.Fatalf("reconciler calls = %d; want 0", reconciler.calls)
	}
}

func TestReconcileOnceRejectsNilSourceAndReconciler(t *testing.T) {
	reconciler := &recordingReconciler{}
	source := &fakeSource{snapshot: readySourceSnapshot(t, 1)}
	request := testRequest(1)

	requireErrorIs(t, ReconcileOnce(context.Background(), request, nil, reconciler).Err, ErrNilSource)
	requireErrorIs(t, ReconcileOnce(context.Background(), request, source, nil).Err, ErrNilReconciler)
	if source.reads != 0 {
		t.Fatalf("source reads = %d; want 0", source.reads)
	}
	if reconciler.calls != 0 {
		t.Fatalf("reconciler calls = %d; want 0", reconciler.calls)
	}
}

func TestReconcileOnceWithNilReconcileFuncReturnsNilReconcilerError(t *testing.T) {
	var fn ReconcileFunc
	source := &fakeSource{snapshot: readySourceSnapshot(t, 1)}

	result := ReconcileOnce(context.Background(), testRequest(1), source, fn)

	requireErrorIs(t, result.Err, ErrNilReconciler)
	if source.reads != 0 {
		t.Fatalf("source reads = %d; want 0", source.reads)
	}
}

func TestReconcileOncePreservesSourceErrorAndSkipsReconcile(t *testing.T) {
	readErr := errors.New("read failed")
	source := &fakeSource{err: readErr}
	reconciler := &recordingReconciler{}

	result := ReconcileOnce(context.Background(), testRequest(1), source, reconciler)

	if result.Err != readErr {
		t.Fatalf("Err = %v; want %v", result.Err, readErr)
	}
	if reconciler.calls != 0 {
		t.Fatalf("reconciler calls = %d; want 0", reconciler.calls)
	}
}

func TestReconcileOncePreservesObjectCacheNotReady(t *testing.T) {
	cache, err := objectcache.New(testCollection())
	requireNoError(t, err)
	reconciler := &recordingReconciler{}

	result := ReconcileOnce(context.Background(), testRequest(1), cache, reconciler)

	requireErrorIs(t, result.Err, objectcache.ErrNotReady)
	if reconciler.calls != 0 {
		t.Fatalf("reconciler calls = %d; want 0", reconciler.calls)
	}
}

func TestReconcileOnceCallsReconcilerOnceWithSourceSnapshot(t *testing.T) {
	sourceSnapshot := readySourceSnapshot(t, 9)
	source := &fakeSource{snapshot: sourceSnapshot}
	reconciler := &recordingReconciler{result: Success()}
	request := testRequest(7)

	result := ReconcileOnce(context.Background(), request, source, reconciler)

	if result.Failed() {
		t.Fatalf("result = %#v; want success", result)
	}
	if source.reads != 1 {
		t.Fatalf("source reads = %d; want 1", source.reads)
	}
	if reconciler.calls != 1 {
		t.Fatalf("reconciler calls = %d; want 1", reconciler.calls)
	}
	if !reconciler.request.Key.Equal(request.Key) {
		t.Fatalf("reconciler request = %#v; want %#v", reconciler.request, request)
	}
	if reconciler.snap.Revision != sourceSnapshot.Revision {
		t.Fatalf("reconciler revision = %s; want %s", reconciler.snap.Revision, sourceSnapshot.Revision)
	}
	if reconciler.snap.View.Revision() != sourceSnapshot.Value.Revision() {
		t.Fatalf("reconciler view revision = %s; want %s", reconciler.snap.View.Revision(), sourceSnapshot.Value.Revision())
	}
}

func TestReconcileOnceCancelledAfterReadSkipsReconcile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	source := &fakeSource{
		snapshot: readySourceSnapshot(t, 3),
		afterRead: func() {
			cancel()
		},
	}
	reconciler := &recordingReconciler{}

	result := ReconcileOnce(ctx, testRequest(1), source, reconciler)

	requireErrorIs(t, result.Err, context.Canceled)
	if source.reads != 1 {
		t.Fatalf("source reads = %d; want 1", source.reads)
	}
	if reconciler.calls != 0 {
		t.Fatalf("reconciler calls = %d; want 0", reconciler.calls)
	}
}

func TestReconcileOnceReturnsUserResultAsIs(t *testing.T) {
	userErr := errors.New("user failed")
	source := &fakeSource{snapshot: readySourceSnapshot(t, 2)}
	reconciler := &recordingReconciler{result: Failure(userErr)}

	result := ReconcileOnce(context.Background(), testRequest(1), source, reconciler)

	if result.Err != userErr {
		t.Fatalf("Err = %v; want %v", result.Err, userErr)
	}
	requireErrorIs(t, result.Err, userErr)
}

func TestReconcileOncePropagatesPanic(t *testing.T) {
	defer func() {
		if got := recover(); got != "boom" {
			t.Fatalf("recover() = %#v; want boom", got)
		}
	}()

	_ = ReconcileOnce(
		context.Background(),
		testRequest(1),
		&fakeSource{snapshot: readySourceSnapshot(t, 1)},
		ReconcileFunc(func(context.Context, Request, Snapshot) Result {
			panic("boom")
		}),
	)
}

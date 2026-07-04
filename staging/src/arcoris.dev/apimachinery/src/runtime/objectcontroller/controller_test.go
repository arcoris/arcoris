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
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/snapshot"
)

func TestControllerStoresConstructionDependencies(t *testing.T) {
	queue := &recordingQueue{}
	source := &fakeSnapshotSource{snapshot: testSnapshot(t, 1)}
	reconciler := &fakeReconciler{}
	controller, err := New(
		Options{Workers: 2},
		queue,
		source,
		reconciler,
	)
	requireNoError(t, err)

	if controller.queue != queue || controller.source != source || controller.reconciler != reconciler || controller.workers != 2 {
		t.Fatalf("controller fields were not initialized from constructor inputs")
	}
}

type fakeSnapshotSource struct {
	mu       sync.Mutex
	snapshot snapshot.Snapshot[objectstore.Revision, objectcache.View]
	err      error
	calls    int
}

func (s *fakeSnapshotSource) ReadSnapshot() (snapshot.Snapshot[objectstore.Revision, objectcache.View], error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.calls++
	if s.err != nil {
		return snapshot.Snapshot[objectstore.Revision, objectcache.View]{}, s.err
	}
	return s.snapshot, nil
}

func (s *fakeSnapshotSource) callCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.calls
}

type fakeReconciler struct {
	mu         sync.Mutex
	result     objectreconciler.Result
	panicValue any
	calls      int
	requests   []objectreconciler.Request
	snaps      []objectreconciler.Snapshot
	started    chan struct{}
	block      <-chan struct{}
}

func (r *fakeReconciler) Reconcile(
	ctx context.Context,
	request objectreconciler.Request,
	snap objectreconciler.Snapshot,
) objectreconciler.Result {
	r.mu.Lock()
	r.calls++
	r.requests = append(r.requests, request)
	r.snaps = append(r.snaps, snap)
	started := r.started
	block := r.block
	panicValue := r.panicValue
	result := r.result
	r.mu.Unlock()

	if started != nil {
		select {
		case started <- struct{}{}:
		default:
		}
	}
	if block != nil {
		select {
		case <-block:
		case <-ctx.Done():
			return objectreconciler.Failure(ctx.Err())
		}
	}
	if panicValue != nil {
		panic(panicValue)
	}
	return result
}

func (r *fakeReconciler) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.calls
}

func (r *fakeReconciler) reconcilerRequests() []objectreconciler.Request {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]objectreconciler.Request(nil), r.requests...)
}

func (r *fakeReconciler) snapshots() []objectreconciler.Snapshot {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]objectreconciler.Snapshot(nil), r.snaps...)
}

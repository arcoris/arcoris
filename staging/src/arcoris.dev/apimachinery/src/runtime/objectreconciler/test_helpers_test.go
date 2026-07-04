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
	"fmt"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/snapshot"
)

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
}

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		Scope: objectstore.AllNamespaces(),
	}
}

func testKey(id int) objectstore.Key {
	return objectstore.MustKey(testCollection().Resource, metaidentity.ObjectName{
		Namespace: "default",
		Name:      metaidentity.Name(fmt.Sprintf("worker-%d", id)),
	})
}

func testRequest(id int) Request {
	return Request{Key: testKey(id)}
}

func readyCache(t testing.TB, revision objectstore.Revision) *objectcache.Cache {
	t.Helper()

	cache, err := objectcache.New(testCollection())
	requireNoError(t, err)
	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, revision)))

	return cache
}

func collectionRead(t testing.TB, revision objectstore.Revision) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(testCollection(), objectstore.ListResult{Revision: revision})
	requireNoError(t, err)

	return read
}

func readySourceSnapshot(t testing.TB, revision objectstore.Revision) snapshot.Snapshot[objectstore.Revision, objectcache.View] {
	t.Helper()

	snap, err := readyCache(t, revision).ReadSnapshot()
	requireNoError(t, err)

	return snap
}

type fakeSource struct {
	snapshot  snapshot.Snapshot[objectstore.Revision, objectcache.View]
	err       error
	afterRead func()
	reads     int
}

func (s *fakeSource) ReadSnapshot() (snapshot.Snapshot[objectstore.Revision, objectcache.View], error) {
	s.reads++
	if s.afterRead != nil {
		s.afterRead()
	}
	if s.err != nil {
		return snapshot.Snapshot[objectstore.Revision, objectcache.View]{}, s.err
	}
	return s.snapshot, nil
}

type recordingReconciler struct {
	result  Result
	calls   int
	ctx     context.Context
	request Request
	snap    Snapshot
}

func (r *recordingReconciler) Reconcile(ctx context.Context, request Request, snap Snapshot) Result {
	r.calls++
	r.ctx = ctx
	r.request = request
	r.snap = snap
	return r.result
}

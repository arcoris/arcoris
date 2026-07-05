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

package objectcontrollerwiring_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectcontrollerwiring"
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestSameObjectSmokeUsesCacheBeforeEnqueueAndPredicate(t *testing.T) {
	keyA := wiringKey("worker-a")
	keyB := wiringKey("worker-b")
	watchErr := errors.New("watch opened")
	reconciler := &recordingReconciler{}

	config := objectcontrollerwiring.SameObjectConfig{
		Source: &scriptedListerWatcher{
			listReads: []storewatchapi.CollectionRead{
				wiringRead(t, 10,
					wiringItem(keyA, 1, "a"),
					wiringItem(keyB, 2, "b"),
				),
			},
			watches: []watchScript{
				{err: watchErr},
			},
		},
		Collection: wiringCollection(),
		Reconciler: reconciler,
		Queue: objectworkqueue.Options{
			Capacity: 4,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
		Predicate: keyPredicate(t, keyA),
	}
	wiring, err := objectcontrollerwiring.NewSameObject(config)
	requireNoError(t, err)

	requireErrorIs(t, wiring.Reflector().Run(context.Background()), watchErr)

	wiring.Queue().ShutDown()
	requireNoError(t, wiring.Controller().Run(context.Background()))

	records := reconciler.recorded()
	requireRecordKeys(t, records, keyA)
	requireSnapshotContains(t, records[0].snapshot, keyA, 1)
	requireSnapshotContains(t, records[0].snapshot, keyB, 2)

	stats := wiring.Queue().Stats()
	if stats.Queued != 0 || stats.Processing != 0 {
		t.Fatalf("queue stats = %#v; want drained queue", stats)
	}
}

func TestSameObjectIndexesObserveReplaceBeforeWorkIsQueued(t *testing.T) {
	key := wiringKey("worker-a")
	watchErr := errors.New("watch opened")
	index := newWiringDesiredIndex(t)
	config := objectcontrollerwiring.SameObjectConfig{
		Source: &scriptedListerWatcher{
			listReads: []storewatchapi.CollectionRead{
				wiringRead(t, 10, wiringItem(key, 1, "group-a")),
			},
			watches: []watchScript{
				{err: watchErr},
			},
		},
		Collection: wiringCollection(),
		Reconciler: &recordingReconciler{},
		Queue: objectworkqueue.Options{
			Capacity: 4,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
		Indexes: []*objectindex.Index{index},
	}
	wiring, err := objectcontrollerwiring.NewSameObject(config)
	requireNoError(t, err)

	requireErrorIs(t, wiring.Reflector().Run(context.Background()), watchErr)

	requireWiringIndexKeys(t, index, "desired", "group-a", key)
	requireSnapshotContains(t, readSameObjectCacheSnapshot(t, wiring), key, 1)
	item, err := wiring.Queue().Get(context.Background())
	requireNoError(t, err)
	if !item.Key.Equal(key) {
		t.Fatalf("queued key = %#v; want %#v", item.Key, key)
	}
	requireNoError(t, wiring.Queue().Done(item))
}

type reconciliationRecord struct {
	request  objectreconciler.Request
	snapshot objectreconciler.Snapshot
}

type recordingReconciler struct {
	mu      sync.Mutex
	records []reconciliationRecord
}

func (r *recordingReconciler) Reconcile(
	_ context.Context,
	request objectreconciler.Request,
	snapshot objectreconciler.Snapshot,
) objectreconciler.Result {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.records = append(r.records, reconciliationRecord{
		request:  request,
		snapshot: snapshot,
	})

	return objectreconciler.Success()
}

func (r *recordingReconciler) recorded() []reconciliationRecord {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]reconciliationRecord(nil), r.records...)
}

type scriptedListerWatcher struct {
	mu sync.Mutex

	listReads []storewatchapi.CollectionRead
	watches   []watchScript
}

func (s *scriptedListerWatcher) ListCollection(
	_ context.Context,
	_ objectstore.ListRequest,
) (storewatchapi.CollectionRead, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.listReads) == 0 {
		return storewatchapi.CollectionRead{}, errors.New("unexpected ListCollection call")
	}
	read := s.listReads[0]
	s.listReads = s.listReads[1:]

	return read, nil
}

func (s *scriptedListerWatcher) Watch(
	_ context.Context,
	_ objectwatch.Request,
) (objectwatch.Stream, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.watches) == 0 {
		return nil, errors.New("unexpected Watch call")
	}
	watch := s.watches[0]
	s.watches = s.watches[1:]

	return nil, watch.err
}

type watchScript struct {
	err error
}

func requireRecordKeys(t testing.TB, records []reconciliationRecord, want ...objectstore.Key) {
	t.Helper()

	if len(records) != len(want) {
		t.Fatalf("reconcile records = %d; want %d", len(records), len(want))
	}
	for i, key := range want {
		if !records[i].request.Key.Equal(key) {
			t.Fatalf("record %d request key = %#v; want %#v", i, records[i].request.Key, key)
		}
	}
}

func requireSnapshotContains(
	t testing.TB,
	snapshot objectreconciler.Snapshot,
	key objectstore.Key,
	revision objectstore.Revision,
) {
	t.Helper()

	result, err := snapshot.View.Get(key)
	requireNoError(t, err)
	if !result.Found {
		t.Fatalf("snapshot revision %s does not contain %#v", snapshot.Revision, key)
	}
	if result.State.Revision != revision {
		t.Fatalf("snapshot state revision = %s; want %s", result.State.Revision, revision)
	}
}

func readSameObjectCacheSnapshot(
	t testing.TB,
	wiring *objectcontrollerwiring.SameObject,
) objectreconciler.Snapshot {
	t.Helper()

	snapshot, err := wiring.Cache().ReadSnapshot()
	requireNoError(t, err)

	return objectreconciler.Snapshot{
		View:     snapshot.Value,
		Revision: snapshot.Revision,
	}
}

func newWiringDesiredIndex(t testing.TB) *objectindex.Index {
	t.Helper()

	index, err := objectindex.New(objectindex.Definition{
		Name: "desired",
		Extractor: objectindex.ExtractorFunc(
			func(item objectstore.ListItem, emit objectindex.EmitFunc) error {
				desired, ok := item.State.Object.Desired.AsString()
				if !ok {
					return errors.New("desired is not string")
				}

				return emit(objectindex.Value(desired))
			},
		),
	})
	requireNoError(t, err)

	return index
}

func requireWiringIndexKeys(
	t testing.TB,
	index *objectindex.Index,
	name objectindex.Name,
	value objectindex.Value,
	want ...objectstore.Key,
) {
	t.Helper()

	keys, err := index.Lookup(name, value)
	requireNoError(t, err)
	if len(keys) != len(want) {
		t.Fatalf("index keys = %d; want %d: %#v", len(keys), len(want), keys)
	}
	for i, key := range want {
		if !keys[i].Equal(key) {
			t.Fatalf("index key %d = %#v; want %#v", i, keys[i], key)
		}
	}
}

func wiringCollection() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: wiringResource(),
		Scope:    objectstore.AllNamespaces(),
	}
}

func wiringResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

func wiringKey(name string) objectstore.Key {
	return objectstore.MustKey(wiringResource(), metaidentity.ObjectName{
		Namespace: "default",
		Name:      metaidentity.Name(name),
	})
}

func wiringItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{
		Key:   key,
		State: wiringState(key, revision, desired),
	}
}

func wiringState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    "Worker",
			}),
			meta.ObjectMeta{
				Name:      key.Object.Name,
				Namespace: key.Object.Namespace,
			},
			value.StringValue(desired),
			value.StringValue(fmt.Sprintf("observed-%s", desired)),
		),
		Revision: revision,
	}
}

func wiringRead(
	t testing.TB,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(wiringCollection(), objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func keyPredicate(t testing.TB, key objectstore.Key) objectquery.Predicate {
	t.Helper()

	query, err := objectquery.KeyEquals(key)
	requireNoError(t, err)
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	return predicate
}

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

package objectcontroller_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectreflectorsink"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestObjectControllerWiringInitialList(t *testing.T) {
	collection := wiringCollection()
	keyA := wiringKey("worker-a")
	keyB := wiringKey("worker-b")

	cache, queue, fanout := newWiringPipeline(t, collection, 4)
	watchOpened := make(chan struct{})
	source := &scriptedListerWatcher{
		listReads: []storewatchapi.CollectionRead{
			wiringRead(t, collection, 10,
				wiringItem(keyA, 1, "a"),
				wiringItem(keyB, 2, "b"),
			),
		},
		watches: []watchScript{
			{stream: &blockingStream{}, opened: watchOpened},
		},
	}
	reflector := newWiringReflector(t, source, collection, fanout)

	cancel, done := startReflector(t, reflector)
	waitSignal(t, watchOpened)
	stopReflector(t, cancel, done)

	queue.ShutDown()
	records := runWiringController(t, queue, cache)

	requireRecordKeys(t, records, keyA, keyB)
	requireSnapshotContains(t, records[0].snapshot, keyA, 1)
	requireSnapshotContains(t, records[0].snapshot, keyB, 2)
}

func TestObjectControllerWiringRelistMissingKeyRepair(t *testing.T) {
	collection := wiringCollection()
	keyA := wiringKey("worker-a")
	keyB := wiringKey("worker-b")
	keyC := wiringKey("worker-c")

	cache, queue, fanout := newWiringPipeline(t, collection, 8)
	firstWatchOpened := make(chan struct{})
	releaseRestart := make(chan struct{})
	secondWatchOpened := make(chan struct{})
	source := &scriptedListerWatcher{
		listReads: []storewatchapi.CollectionRead{
			wiringRead(t, collection, 10,
				wiringItem(keyA, 1, "a"),
				wiringItem(keyB, 2, "b"),
			),
			wiringRead(t, collection, 20,
				wiringItem(keyB, 12, "b2"),
				wiringItem(keyC, 13, "c"),
			),
		},
		watches: []watchScript{
			{
				stream: &gatedEventStream{
					release: releaseRestart,
					event:   wiringRestartEvent(t),
				},
				opened: firstWatchOpened,
			},
			{stream: &blockingStream{}, opened: secondWatchOpened},
		},
	}
	reflector := newWiringReflector(t, source, collection, fanout)

	cancel, done := startReflector(t, reflector)
	waitSignal(t, firstWatchOpened)
	drainQueueItems(t, queue, keyA, keyB)
	close(releaseRestart)
	waitSignal(t, secondWatchOpened)
	stopReflector(t, cancel, done)

	queue.ShutDown()
	records := runWiringController(t, queue, cache)

	requireRecordKeys(t, records, keyB, keyC, keyA)
	requireSnapshotContains(t, records[2].snapshot, keyB, 12)
	requireSnapshotContains(t, records[2].snapshot, keyC, 13)
	requireSnapshotMissing(t, records[2].snapshot, keyA)
}

func TestObjectControllerWiringApplyChange(t *testing.T) {
	collection := wiringCollection()
	keyA := wiringKey("worker-a")
	before := wiringState(keyA, 1, "a")
	after := wiringState(keyA, 2, "a2")
	change := objectstore.MustUpdatedChange(keyA, before, after)

	cache, queue, fanout := newWiringPipeline(t, collection, 4)
	watchOpened := make(chan struct{})
	releaseChange := make(chan struct{})
	changeApplied := make(chan struct{})
	source := &scriptedListerWatcher{
		listReads: []storewatchapi.CollectionRead{
			wiringRead(t, collection, 1, objectstore.ListItem{Key: keyA, State: before}),
		},
		watches: []watchScript{
			{
				stream: &gatedEventStream{
					release: releaseChange,
					event:   wiringChangedEvent(t, change),
					after:   changeApplied,
				},
				opened: watchOpened,
			},
		},
	}
	reflector := newWiringReflector(t, source, collection, fanout)

	cancel, done := startReflector(t, reflector)
	waitSignal(t, watchOpened)
	drainQueueItems(t, queue, keyA)
	close(releaseChange)
	waitSignal(t, changeApplied)
	stopReflector(t, cancel, done)

	queue.ShutDown()
	records := runWiringController(t, queue, cache)

	requireRecordKeys(t, records, keyA)
	requireSnapshotContains(t, records[0].snapshot, keyA, 2)
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
	signal(watch.opened)

	return watch.stream, nil
}

type watchScript struct {
	stream objectwatch.Stream
	opened chan<- struct{}
}

type blockingStream struct{}

func (*blockingStream) Next(ctx context.Context) (objectwatch.Event, error) {
	<-ctx.Done()
	return objectwatch.Event{}, ctx.Err()
}

func (*blockingStream) Close() error {
	return nil
}

type gatedEventStream struct {
	mu       sync.Mutex
	release  <-chan struct{}
	event    objectwatch.Event
	after    chan<- struct{}
	returned bool
}

func (s *gatedEventStream) Next(ctx context.Context) (objectwatch.Event, error) {
	s.mu.Lock()
	if !s.returned {
		s.returned = true
		release := s.release
		event := s.event
		s.mu.Unlock()

		select {
		case <-release:
			return event, nil
		case <-ctx.Done():
			return objectwatch.Event{}, ctx.Err()
		}
	}
	after := s.after
	s.after = nil
	s.mu.Unlock()

	signal(after)
	<-ctx.Done()
	return objectwatch.Event{}, ctx.Err()
}

func (*gatedEventStream) Close() error {
	return nil
}

func newWiringPipeline(
	t testing.TB,
	collection objectstore.ListRequest,
	capacity int,
) (*objectcache.Cache, *objectworkqueue.Queue, *objectreflectorsink.Fanout) {
	t.Helper()

	cache, err := objectcache.New(collection)
	requireNoError(t, err)
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: capacity})
	requireNoError(t, err)
	enqueueSink, err := objectenqueue.NewReflectorSink(objectenqueue.ReflectorSinkConfig{
		Queue:   queue,
		Listed:  objectenqueue.ListedObject(),
		Changed: objectenqueue.ChangedObject(),
	})
	requireNoError(t, err)
	fanout, err := objectreflectorsink.NewFanout(cache, enqueueSink)
	requireNoError(t, err)

	return cache, queue, fanout
}

func newWiringReflector(
	t testing.TB,
	source storewatchapi.ListerWatcher,
	collection objectstore.ListRequest,
	sink *objectreflectorsink.Fanout,
) *objectreflector.Reflector {
	t.Helper()

	reflector, err := objectreflector.New(source, collection, sink)
	requireNoError(t, err)

	return reflector
}

func startReflector(
	t testing.TB,
	reflector *objectreflector.Reflector,
) (context.CancelFunc, <-chan error) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- reflector.Run(ctx)
	}()

	return cancel, done
}

func stopReflector(t testing.TB, cancel context.CancelFunc, done <-chan error) {
	t.Helper()

	cancel()
	err := waitError(t, done)
	requireErrorIs(t, err, context.Canceled)
}

func runWiringController(
	t testing.TB,
	queue *objectworkqueue.Queue,
	source objectreconciler.SnapshotSource,
) []reconciliationRecord {
	t.Helper()

	reconciler := &recordingReconciler{}
	controller, err := objectcontroller.New(
		objectcontroller.Options{Workers: 1},
		queue,
		source,
		reconciler,
	)
	requireNoError(t, err)
	requireNoError(t, controller.Run(context.Background()))

	stats := queue.Stats()
	if stats.Queued != 0 || stats.Processing != 0 {
		t.Fatalf("queue stats = %#v; want drained queue", stats)
	}

	return reconciler.recorded()
}

func drainQueueItems(t testing.TB, queue *objectworkqueue.Queue, want ...objectstore.Key) {
	t.Helper()

	for _, key := range want {
		item, err := queue.Get(context.Background())
		requireNoError(t, err)
		if !item.Key.Equal(key) {
			t.Fatalf("queue item key = %#v; want %#v", item.Key, key)
		}
		requireNoError(t, queue.Done(item))
	}
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

func requireSnapshotMissing(t testing.TB, snapshot objectreconciler.Snapshot, key objectstore.Key) {
	t.Helper()

	result, err := snapshot.View.Get(key)
	requireNoError(t, err)
	if result.Found {
		t.Fatalf("snapshot revision %s unexpectedly contains %#v", snapshot.Revision, key)
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
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}

func wiringRead(
	t testing.TB,
	collection objectstore.ListRequest,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(collection, objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func wiringChangedEvent(t testing.TB, change objectstore.Change) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Changed(change)
	requireNoError(t, err)

	return event
}

func wiringRestartEvent(t testing.TB) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.RestartRequired(objectwatch.RestartContinuityLost, 0)
	requireNoError(t, err)

	return event
}

func waitSignal(t testing.TB, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for signal")
	}
}

func waitError(t testing.TB, ch <-chan error) error {
	t.Helper()

	select {
	case err := <-ch:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for error")
		return nil
	}
}

func signal(ch chan<- struct{}) {
	if ch == nil {
		return
	}
	close(ch)
}

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

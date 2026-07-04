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

package objectcontrollerwiring

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
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestRunSameObjectPanicsOnNilContext(t *testing.T) {
	graph := newRunTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})

	defer func() {
		if recover() == nil {
			t.Fatal("RunSameObject did not panic")
		}
	}()

	_ = RunSameObject(nil, graph)
}

func TestRunSameObjectRejectsInvalidGraph(t *testing.T) {
	tests := []struct {
		name  string
		graph func(t testing.TB) *SameObject
	}{
		{
			name: "nil graph",
			graph: func(testing.TB) *SameObject {
				return nil
			},
		},
		{
			name: "nil queue",
			graph: func(t testing.TB) *SameObject {
				graph := newRunTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})
				graph.queue = nil
				return graph
			},
		},
		{
			name: "nil reflector",
			graph: func(t testing.TB) *SameObject {
				graph := newRunTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})
				graph.reflector = nil
				return graph
			},
		},
		{
			name: "nil controller",
			graph: func(t testing.TB) *SameObject {
				graph := newRunTestGraph(t, &runTestListerWatcher{}, &runTestReconciler{})
				graph.controller = nil
				return graph
			},
		},
	}

	for _, tt := range tests {
		err := RunSameObject(context.Background(), tt.graph(t))

		if !errors.Is(err, ErrInvalidSameObject) {
			t.Fatalf("%s: error = %v; want errors.Is(%v)", tt.name, err, ErrInvalidSameObject)
		}
	}
}

func TestRunSameObjectStartsReflectorAndController(t *testing.T) {
	key := runTestKey("alpha")
	reconciler := &runTestReconciler{
		started: make(chan struct{}),
	}
	source := &runTestListerWatcher{
		read:        runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newRunTestGraph(t, source, reconciler)
	ctx, cancel := context.WithCancel(context.Background())

	result := runSameObjectAsync(ctx, graph)
	waitForSignal(t, source.watchCalled)
	waitForSignal(t, reconciler.started)
	cancel()

	requireRunResult(t, result, context.Canceled)
	requireReconcilerKeys(t, reconciler, key)
	requireListWatchCalls(t, source, 1, 1)
}

func TestRunSameObjectShutsDownQueueWhenReflectorReturns(t *testing.T) {
	key := runTestKey("alpha")
	reflectorErr := errors.New("reflector failed")
	reconciler := &runTestReconciler{
		started: make(chan struct{}),
	}
	source := &runTestListerWatcher{
		read:            runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		watchWaitStrict: reconciler.started,
		watchErr:        reflectorErr,
	}
	graph := newRunTestGraph(t, source, reconciler)

	err := RunSameObject(context.Background(), graph)

	requireErrorIs(t, err, reflectorErr)
	if !graph.Queue().IsShutDown() {
		t.Fatal("queue is not shut down")
	}
	requireReconcilerKeys(t, reconciler, key)
	requireSnapshotContains(t, reconciler, key, 1)
	requireQueueDrained(t, graph.Queue())
}

func TestRunSameObjectCancelsReflectorWhenControllerReturnsFirst(t *testing.T) {
	key := runTestKey("alpha")
	controllerErr := errors.New("controller failed")
	stream := runTestWaitingStream()
	reconciler := &runTestReconciler{
		err: controllerErr,
	}
	source := &runTestListerWatcher{
		read:   runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		stream: stream,
	}
	graph := newRunTestGraph(t, source, reconciler)

	err := RunSameObject(context.Background(), graph)

	requireErrorIs(t, err, controllerErr)
	waitForSignal(t, stream.done)
	if !graph.Queue().IsShutDown() {
		t.Fatal("queue is not shut down")
	}
	requireReconcilerKeys(t, reconciler, key)
}

func TestRunSameObjectWaitsForBothSides(t *testing.T) {
	key := runTestKey("alpha")
	controllerErr := errors.New("controller failed")
	releaseStream := make(chan struct{})
	stream := &runTestHeldStream{
		entered: make(chan struct{}),
		done:    make(chan struct{}),
		release: releaseStream,
	}
	reconciler := &runTestReconciler{
		err:      controllerErr,
		returned: make(chan struct{}),
	}
	source := &runTestListerWatcher{
		read:   runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		stream: stream,
	}
	graph := newRunTestGraph(t, source, reconciler)

	result := runSameObjectAsync(context.Background(), graph)
	waitForSignal(t, stream.entered)
	waitForSignal(t, reconciler.returned)
	requireNoRunResultYet(t, result)

	close(releaseStream)
	requireRunResult(t, result, controllerErr)
	waitForSignal(t, stream.done)
}

func TestRunSameObjectReturnsParentContextError(t *testing.T) {
	source := &runTestListerWatcher{
		read:        runTestRead(t, 10),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newRunTestGraph(t, source, &runTestReconciler{})
	ctx, cancel := context.WithCancel(context.Background())

	result := runSameObjectAsync(ctx, graph)
	waitForSignal(t, source.watchCalled)
	cancel()

	requireRunResult(t, result, context.Canceled)
}

func TestRunSameObjectReturnsReflectorFatalError(t *testing.T) {
	reflectorErr := errors.New("reflector failed")
	source := &runTestListerWatcher{
		read:     runTestRead(t, 10),
		watchErr: reflectorErr,
	}
	graph := newRunTestGraph(t, source, &runTestReconciler{})

	err := RunSameObject(context.Background(), graph)

	requireErrorIs(t, err, reflectorErr)
}

func TestRunSameObjectReturnsControllerFatalError(t *testing.T) {
	key := runTestKey("alpha")
	controllerErr := errors.New("controller failed")
	source := &runTestListerWatcher{
		read:   runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		stream: runTestWaitingStream(),
	}
	graph := newRunTestGraph(t, source, &runTestReconciler{err: controllerErr})

	err := RunSameObject(context.Background(), graph)

	requireErrorIs(t, err, controllerErr)
}

func TestRunSameObjectJoinsFatalErrors(t *testing.T) {
	key := runTestKey("alpha")
	reflectorErr := errors.New("reflector failed")
	controllerErr := errors.New("controller failed")
	reconciler := &runTestReconciler{
		err:            controllerErr,
		started:        make(chan struct{}),
		waitForContext: true,
	}
	source := &runTestListerWatcher{
		read:            runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		watchWaitStrict: reconciler.started,
		watchErr:        reflectorErr,
	}
	graph := newRunTestGraph(t, source, reconciler)

	result := runSameObjectAsync(context.Background(), graph)
	waitForSignal(t, reconciler.started)

	err := readRunResult(t, result)
	requireErrorIs(t, err, reflectorErr)
	requireErrorIs(t, err, controllerErr)
	requireReconcilerKeys(t, reconciler, key)
}

func TestRunSameObjectStartsEachSideOnce(t *testing.T) {
	key := runTestKey("alpha")
	reflectorErr := errors.New("reflector failed")
	reconciler := &runTestReconciler{
		started: make(chan struct{}),
	}
	source := &runTestListerWatcher{
		read:            runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		watchWaitStrict: reconciler.started,
		watchErr:        reflectorErr,
	}
	graph := newRunTestGraph(t, source, reconciler)

	requireErrorIs(t, RunSameObject(context.Background(), graph), reflectorErr)

	requireListWatchCalls(t, source, 1, 1)
	if calls := reconciler.callCount(); calls != 1 {
		t.Fatalf("reconciler calls = %d; want 1", calls)
	}
}

func TestRunSameObjectSmokeWithNewSameObject(t *testing.T) {
	key := runTestKey("alpha")
	reconciler := &runTestReconciler{
		started: make(chan struct{}),
	}
	source := &runTestListerWatcher{
		read:        runTestRead(t, 10, runTestItem(key, 1, "alpha")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newRunTestGraph(t, source, reconciler)
	ctx, cancel := context.WithCancel(context.Background())

	result := runSameObjectAsync(ctx, graph)
	waitForSignal(t, source.watchCalled)
	waitForSignal(t, reconciler.started)
	cancel()

	requireRunResult(t, result, context.Canceled)
	requireReconcilerKeys(t, reconciler, key)
	requireSnapshotContains(t, reconciler, key, 1)
}

type runTestResult struct {
	err error
}

func runSameObjectAsync(ctx context.Context, graph *SameObject) <-chan runTestResult {
	result := make(chan runTestResult, 1)
	go func() {
		result <- runTestResult{err: RunSameObject(ctx, graph)}
	}()

	return result
}

func requireRunResult(t testing.TB, result <-chan runTestResult, target error) {
	t.Helper()

	requireErrorIs(t, readRunResult(t, result), target)
}

func readRunResult(t testing.TB, result <-chan runTestResult) error {
	t.Helper()

	select {
	case result := <-result:
		return result.err
	case <-time.After(5 * time.Second):
		t.Fatal("RunSameObject did not return")
		return nil
	}
}

func requireNoRunResultYet(t testing.TB, result <-chan runTestResult) {
	t.Helper()

	select {
	case result := <-result:
		t.Fatalf("RunSameObject returned early with %v", result.err)
	default:
	}
}

func waitForSignal(t testing.TB, ch <-chan struct{}) {
	t.Helper()

	select {
	case <-ch:
	case <-time.After(5 * time.Second):
		t.Fatal("signal was not received")
	}
}

func newRunTestGraph(
	t testing.TB,
	source storewatchapi.ListerWatcher,
	reconciler objectreconciler.Reconciler,
) *SameObject {
	t.Helper()

	graph, err := NewSameObject(SameObjectConfig{
		Source:     source,
		Collection: runTestCollection(),
		Reconciler: reconciler,
		Queue: objectworkqueue.Options{
			Capacity: 8,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
	})
	requireNoError(t, err)

	return graph
}

type runTestListerWatcher struct {
	mu sync.Mutex

	read            storewatchapi.CollectionRead
	listErr         error
	stream          objectwatch.Stream
	watchWait       <-chan struct{}
	watchWaitStrict <-chan struct{}
	watchErr        error

	listCalls  int
	watchCalls int

	watchCalled chan struct{}
	watchOnce   sync.Once
}

func (s *runTestListerWatcher) ListCollection(
	_ context.Context,
	_ objectstore.ListRequest,
) (storewatchapi.CollectionRead, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.listCalls++
	if s.listErr != nil {
		return storewatchapi.CollectionRead{}, s.listErr
	}

	return s.read, nil
}

func (s *runTestListerWatcher) Watch(
	ctx context.Context,
	_ objectwatch.Request,
) (objectwatch.Stream, error) {
	s.mu.Lock()
	s.watchCalls++
	watchWait := s.watchWait
	watchWaitStrict := s.watchWaitStrict
	watchErr := s.watchErr
	stream := s.stream
	watchCalled := s.watchCalled
	s.mu.Unlock()

	if watchCalled != nil {
		s.watchOnce.Do(func() { close(watchCalled) })
	}
	if watchWaitStrict != nil {
		<-watchWaitStrict
	} else if watchWait != nil {
		select {
		case <-watchWait:
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if watchErr != nil {
		return nil, watchErr
	}
	if stream != nil {
		return stream, nil
	}

	return runTestWaitingStream(), nil
}

func (s *runTestListerWatcher) callCounts() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.listCalls, s.watchCalls
}

type runTestWaitingWatchStream struct {
	done chan struct{}
	once sync.Once
}

func runTestWaitingStream() *runTestWaitingWatchStream {
	return &runTestWaitingWatchStream{done: make(chan struct{})}
}

func (s *runTestWaitingWatchStream) Next(ctx context.Context) (objectwatch.Event, error) {
	<-ctx.Done()
	s.once.Do(func() { close(s.done) })

	return objectwatch.Event{}, ctx.Err()
}

func (s *runTestWaitingWatchStream) Close() error {
	return nil
}

type runTestHeldStream struct {
	entered chan struct{}
	done    chan struct{}
	release <-chan struct{}
	once    sync.Once
}

func (s *runTestHeldStream) Next(context.Context) (objectwatch.Event, error) {
	close(s.entered)
	<-s.release
	s.once.Do(func() { close(s.done) })

	return objectwatch.Event{}, context.Canceled
}

func (s *runTestHeldStream) Close() error {
	return nil
}

type runTestReconciler struct {
	mu sync.Mutex

	err error

	started        chan struct{}
	startOnce      sync.Once
	returned       chan struct{}
	returnOnce     sync.Once
	release        <-chan struct{}
	waitForContext bool

	records []runTestReconciliation
}

func (r *runTestReconciler) Reconcile(
	ctx context.Context,
	request objectreconciler.Request,
	snapshot objectreconciler.Snapshot,
) objectreconciler.Result {
	if r.started != nil {
		r.startOnce.Do(func() { close(r.started) })
	}
	if r.waitForContext {
		<-ctx.Done()
	}
	if r.release != nil {
		<-r.release
	}

	r.mu.Lock()
	r.records = append(r.records, runTestReconciliation{
		request:  request,
		snapshot: snapshot,
	})
	err := r.err
	r.mu.Unlock()

	if r.returned != nil {
		r.returnOnce.Do(func() { close(r.returned) })
	}
	if err != nil {
		return objectreconciler.Failure(err)
	}

	return objectreconciler.Success()
}

func (r *runTestReconciler) recorded() []runTestReconciliation {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]runTestReconciliation(nil), r.records...)
}

func (r *runTestReconciler) callCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.records)
}

type runTestReconciliation struct {
	request  objectreconciler.Request
	snapshot objectreconciler.Snapshot
}

func requireReconcilerKeys(t testing.TB, reconciler *runTestReconciler, want ...objectstore.Key) {
	t.Helper()

	records := reconciler.recorded()
	if len(records) != len(want) {
		t.Fatalf("reconciler records = %d; want %d", len(records), len(want))
	}
	for i, key := range want {
		if !records[i].request.Key.Equal(key) {
			t.Fatalf("record %d request key = %#v; want %#v", i, records[i].request.Key, key)
		}
	}
}

func requireSnapshotContains(
	t testing.TB,
	reconciler *runTestReconciler,
	key objectstore.Key,
	revision objectstore.Revision,
) {
	t.Helper()

	records := reconciler.recorded()
	if len(records) == 0 {
		t.Fatal("no reconciler records")
	}

	result, err := records[0].snapshot.View.Get(key)
	requireNoError(t, err)
	if !result.Found {
		t.Fatalf("snapshot revision %s does not contain %#v", records[0].snapshot.Revision, key)
	}
	if result.State.Revision != revision {
		t.Fatalf("snapshot state revision = %s; want %s", result.State.Revision, revision)
	}
}

func requireQueueDrained(t testing.TB, queue *objectworkqueue.Queue) {
	t.Helper()

	stats := queue.Stats()
	if stats.Queued != 0 || stats.Processing != 0 {
		t.Fatalf("queue stats = %#v; want drained queue", stats)
	}
}

func requireListWatchCalls(t testing.TB, source *runTestListerWatcher, wantList int, wantWatch int) {
	t.Helper()

	listCalls, watchCalls := source.callCounts()
	if listCalls != wantList || watchCalls != wantWatch {
		t.Fatalf("list/watch calls = %d/%d; want %d/%d", listCalls, watchCalls, wantList, wantWatch)
	}
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

func runTestCollection() objectstore.ListRequest {
	return objectstore.ListRequest{
		Resource: runTestResource(),
		Scope:    objectstore.AllNamespaces(),
	}
}

func runTestResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

func runTestKey(name string) objectstore.Key {
	return objectstore.MustKey(runTestResource(), metaidentity.ObjectName{
		Namespace: "default",
		Name:      metaidentity.Name(name),
	})
}

func runTestItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{
		Key:   key,
		State: runTestState(key, revision, desired),
	}
}

func runTestState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
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

func runTestRead(
	t testing.TB,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(runTestCollection(), objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

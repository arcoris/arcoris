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
	"sync"
	"testing"
	"time"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/runtime/objectcontroller"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestMappedObjectInitialListMapsToMultipleRequestsAndChangeMapsToRequest(t *testing.T) {
	sourceKey := runTestKey("source-a")
	targetX := runTestKey("target-x")
	targetY := runTestKey("target-y")
	targetZ := runTestKey("target-z")
	stream := newMappedWatchStream()
	reconciler := newMappedRecordingReconciler()
	source := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(sourceKey, 1, "source-a")),
		stream:      stream,
		watchCalled: make(chan struct{}),
	}
	graph := newMappedIntegrationGraph(t, source, reconciler,
		listItemMapperForKeys(targetX, targetY),
		changeMapperForKeys(targetZ),
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := runMappedObjectAsync(ctx, graph)
	waitForSignal(t, source.watchCalled)

	records := reconciler.waitForRecords(t, 2)
	requireMappedRecordKeys(t, records[:2], targetX, targetY)
	requireMappedRecordSnapshotContains(t, records[0], sourceKey, 1)
	requireMappedRecordSnapshotContains(t, records[1], sourceKey, 1)

	stream.send(t, mappedChangedEvent(t, updatedMappedChange(t, sourceKey, 1, 2)))
	records = reconciler.waitForRecords(t, 3)
	requireMappedRecordKeys(t, records[2:3], targetZ)
	requireMappedRecordSnapshotContains(t, records[2], sourceKey, 2)

	cancel()
	requireErrorIs(t, readMappedRunResult(t, result), context.Canceled)
}

func TestMappedObjectAllowsZeroMappedWork(t *testing.T) {
	sourceKey := runTestKey("source-a")
	reconciler := newMappedRecordingReconciler()
	source := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(sourceKey, 1, "source-a")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newMappedIntegrationGraph(t, source, reconciler, zeroListItemMapper(), zeroChangeMapper())
	ctx, cancel := context.WithCancel(context.Background())

	result := runMappedObjectAsync(ctx, graph)
	waitForSignal(t, source.watchCalled)
	cancel()

	requireErrorIs(t, readMappedRunResult(t, result), context.Canceled)
	if count := reconciler.recordCount(); count != 0 {
		t.Fatalf("reconciler records = %d; want 0", count)
	}
}

func TestMappedObjectReturnsListedMapperError(t *testing.T) {
	sourceKey := runTestKey("source-a")
	mapperErr := errors.New("listed mapper failed")
	source := &runTestListerWatcher{
		read: runTestRead(t, 1, runTestItem(sourceKey, 1, "source-a")),
	}
	graph := newMappedIntegrationGraph(
		t,
		source,
		newMappedRecordingReconciler(),
		listItemMapperError(mapperErr),
		zeroChangeMapper(),
	)

	err := RunMappedObject(context.Background(), graph)

	requireErrorIs(t, err, mapperErr)
}

func TestMappedObjectReturnsChangedMapperError(t *testing.T) {
	sourceKey := runTestKey("source-a")
	mapperErr := errors.New("changed mapper failed")
	stream := newMappedWatchStream()
	stream.send(t, mappedChangedEvent(t, updatedMappedChange(t, sourceKey, 1, 2)))
	source := &runTestListerWatcher{
		read:   runTestRead(t, 1, runTestItem(sourceKey, 1, "source-a")),
		stream: stream,
	}
	graph := newMappedIntegrationGraph(
		t,
		source,
		newMappedRecordingReconciler(),
		zeroListItemMapper(),
		changeMapperError(mapperErr),
	)

	err := RunMappedObject(context.Background(), graph)

	requireErrorIs(t, err, mapperErr)
}

func newMappedIntegrationGraph(
	t testing.TB,
	source *runTestListerWatcher,
	reconciler objectreconciler.Reconciler,
	listed objectenqueue.ListItemMapper,
	changed objectenqueue.Mapper,
) *MappedObject {
	t.Helper()

	graph, err := NewMappedObject(MappedObjectConfig{
		Source:     source,
		Collection: runTestCollection(),
		Reconciler: reconciler,
		Queue: objectworkqueue.Options{
			Capacity: 8,
		},
		Controller: objectcontroller.Options{
			Workers: 1,
		},
		Listed:  listed,
		Changed: changed,
	})
	requireNoError(t, err)

	return graph
}

// listItemMapperForKeys models source-list mapping without involving query or
// target lookup behavior. Each listed source item produces the same fixed
// target keys so tests can assert queue/controller delivery precisely.
func listItemMapperForKeys(keys ...objectstore.Key) objectenqueue.ListItemMapper {
	return objectenqueue.ListItemMapperFunc(func(_ objectstore.ListItem, emit objectenqueue.EmitFunc) error {
		for _, key := range keys {
			if err := emit(objectworkqueue.Item{Key: key}); err != nil {
				return err
			}
		}

		return nil
	})
}

// changeMapperForKeys is the ApplyChange companion to listItemMapperForKeys.
// It verifies that changed-source mapping travels through the reflector sink,
// real queue, and controller path rather than through direct queue injection.
func changeMapperForKeys(keys ...objectstore.Key) objectenqueue.Mapper {
	return objectenqueue.MapperFunc(func(_ objectstore.Change, emit objectenqueue.EmitFunc) error {
		for _, key := range keys {
			if err := emit(objectworkqueue.Item{Key: key}); err != nil {
				return err
			}
		}

		return nil
	})
}

// zeroListItemMapper proves that "no mapped work" is a valid list-boundary
// result, not an error or an implicit same-object fallback.
func zeroListItemMapper() objectenqueue.ListItemMapper {
	return objectenqueue.ListItemMapperFunc(func(objectstore.ListItem, objectenqueue.EmitFunc) error {
		return nil
	})
}

// zeroChangeMapper proves that an admitted change may intentionally produce no
// mapped target work.
func zeroChangeMapper() objectenqueue.Mapper {
	return objectenqueue.MapperFunc(func(objectstore.Change, objectenqueue.EmitFunc) error {
		return nil
	})
}

// listItemMapperError injects a fatal Replace-path mapping error while keeping
// the source read itself valid.
func listItemMapperError(err error) objectenqueue.ListItemMapper {
	return objectenqueue.ListItemMapperFunc(func(objectstore.ListItem, objectenqueue.EmitFunc) error {
		return err
	})
}

// changeMapperError injects a fatal ApplyChange-path mapping error after cache
// preconditions have been satisfied.
func changeMapperError(err error) objectenqueue.Mapper {
	return objectenqueue.MapperFunc(func(objectstore.Change, objectenqueue.EmitFunc) error {
		return err
	})
}

// mappedRunResult gives asynchronous runner tests a buffered, leak-free handoff
// for the final RunMappedObject error.
type mappedRunResult struct {
	err error
}

// runMappedObjectAsync starts the real runner and lets tests coordinate it with
// watch events or cancellation.
func runMappedObjectAsync(ctx context.Context, graph *MappedObject) <-chan mappedRunResult {
	result := make(chan mappedRunResult, 1)
	go func() {
		result <- mappedRunResult{err: RunMappedObject(ctx, graph)}
	}()

	return result
}

// readMappedRunResult bounds broken-test hangs without making successful tests
// depend on wall-clock timing.
func readMappedRunResult(t testing.TB, result <-chan mappedRunResult) error {
	t.Helper()

	select {
	case result := <-result:
		return result.err
	case <-time.After(5 * time.Second):
		t.Fatal("RunMappedObject did not return")
		return nil
	}
}

// mappedWatchStream is a tiny deterministic watch stream for integration tests.
// Events are pushed by the test, and cancellation exits Next exactly like a
// real reflector watch would.
type mappedWatchStream struct {
	events chan objectwatch.Event
	once   sync.Once
	done   chan struct{}
}

// newMappedWatchStream creates a one-event buffered stream so a test can queue
// the first event before the reflector reaches Next.
func newMappedWatchStream() *mappedWatchStream {
	return &mappedWatchStream{
		events: make(chan objectwatch.Event, 1),
		done:   make(chan struct{}),
	}
}

// send publishes one committed watch event and fails fast if the runner is not
// able to receive test input.
func (s *mappedWatchStream) send(t testing.TB, event objectwatch.Event) {
	t.Helper()

	select {
	case s.events <- event:
	case <-time.After(5 * time.Second):
		t.Fatal("watch stream did not accept event")
	}
}

// Next implements objectwatch.Stream for mappedWatchStream.
func (s *mappedWatchStream) Next(ctx context.Context) (objectwatch.Event, error) {
	select {
	case event := <-s.events:
		return event, nil
	case <-ctx.Done():
		s.once.Do(func() { close(s.done) })
		return objectwatch.Event{}, ctx.Err()
	}
}

// Close implements objectwatch.Stream. The test stream has no external
// resource to release, so Close is intentionally a no-op.
func (s *mappedWatchStream) Close() error {
	return nil
}

// mappedRecordingReconciler records requests and snapshots observed through the
// real controller path. The changed channel advances on every record so tests
// can wait without polling.
type mappedRecordingReconciler struct {
	mu      sync.Mutex
	changed chan struct{}
	records []mappedRecord
}

// newMappedRecordingReconciler prepares an empty recorder with an initial wait
// channel.
func newMappedRecordingReconciler() *mappedRecordingReconciler {
	return &mappedRecordingReconciler{changed: make(chan struct{})}
}

// Reconcile records the exact request and stable snapshot passed by
// objectcontroller, then reports success so the queue can call Done normally.
func (r *mappedRecordingReconciler) Reconcile(
	_ context.Context,
	request objectreconciler.Request,
	snapshot objectreconciler.Snapshot,
) objectreconciler.Result {
	r.mu.Lock()
	r.records = append(r.records, mappedRecord{request: request, snapshot: snapshot})
	close(r.changed)
	r.changed = make(chan struct{})
	r.mu.Unlock()

	return objectreconciler.Success()
}

// waitForRecords returns a copy of all records once at least count records have
// arrived.
func (r *mappedRecordingReconciler) waitForRecords(t testing.TB, count int) []mappedRecord {
	t.Helper()

	deadline := time.After(5 * time.Second)
	for {
		r.mu.Lock()
		if len(r.records) >= count {
			records := append([]mappedRecord(nil), r.records...)
			r.mu.Unlock()
			return records
		}
		changed := r.changed
		r.mu.Unlock()

		select {
		case <-changed:
		case <-deadline:
			t.Fatalf("reconciler records = %d; want at least %d", r.recordCount(), count)
		}
	}
}

// recordCount reports the current number of reconciliation records without
// exposing the recorder's mutable slice.
func (r *mappedRecordingReconciler) recordCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.records)
}

// mappedRecord captures the two values that prove objectcontroller delivered a
// mapped request against the expected source snapshot.
type mappedRecord struct {
	request  objectreconciler.Request
	snapshot objectreconciler.Snapshot
}

// requireMappedRecordKeys asserts mapped target delivery order.
func requireMappedRecordKeys(t testing.TB, records []mappedRecord, want ...objectstore.Key) {
	t.Helper()

	if len(records) != len(want) {
		t.Fatalf("records = %d; want %d", len(records), len(want))
	}
	for i, key := range want {
		if !records[i].request.Key.Equal(key) {
			t.Fatalf("record %d request key = %#v; want %#v", i, records[i].request.Key, key)
		}
	}
}

// requireMappedRecordSnapshotContains asserts cache-before-enqueue behavior by
// checking that the reconciler snapshot already contains the reflected source
// object revision that produced mapped work.
func requireMappedRecordSnapshotContains(
	t testing.TB,
	record mappedRecord,
	key objectstore.Key,
	revision objectstore.Revision,
) {
	t.Helper()

	result, err := record.snapshot.View.Get(key)
	requireNoError(t, err)
	if !result.Found {
		t.Fatalf("snapshot revision %s does not contain %#v", record.snapshot.Revision, key)
	}
	if result.State.Revision != revision {
		t.Fatalf("snapshot state revision = %s; want %s", result.State.Revision, revision)
	}
}

// updatedMappedChange builds a cache-valid update change for sourceKey. The
// before state mirrors the initial list fixture so objectcache accepts the
// transition before the Changed mapper is exercised.
func updatedMappedChange(
	t testing.TB,
	key objectstore.Key,
	beforeRevision objectstore.Revision,
	afterRevision objectstore.Revision,
) objectstore.Change {
	t.Helper()

	return objectstore.MustUpdatedChange(
		key,
		runTestState(key, beforeRevision, string(key.Object.Name)),
		runTestState(key, afterRevision, string(key.Object.Name)+"-updated"),
	)
}

// mappedChangedEvent wraps a committed objectstore.Change in the real watch
// event type expected by objectreflector.
func mappedChangedEvent(t testing.TB, change objectstore.Change) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Changed(change)
	requireNoError(t, err)

	return event
}

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

func zeroListItemMapper() objectenqueue.ListItemMapper {
	return objectenqueue.ListItemMapperFunc(func(objectstore.ListItem, objectenqueue.EmitFunc) error {
		return nil
	})
}

func zeroChangeMapper() objectenqueue.Mapper {
	return objectenqueue.MapperFunc(func(objectstore.Change, objectenqueue.EmitFunc) error {
		return nil
	})
}

func listItemMapperError(err error) objectenqueue.ListItemMapper {
	return objectenqueue.ListItemMapperFunc(func(objectstore.ListItem, objectenqueue.EmitFunc) error {
		return err
	})
}

func changeMapperError(err error) objectenqueue.Mapper {
	return objectenqueue.MapperFunc(func(objectstore.Change, objectenqueue.EmitFunc) error {
		return err
	})
}

type mappedRunResult struct {
	err error
}

func runMappedObjectAsync(ctx context.Context, graph *MappedObject) <-chan mappedRunResult {
	result := make(chan mappedRunResult, 1)
	go func() {
		result <- mappedRunResult{err: RunMappedObject(ctx, graph)}
	}()

	return result
}

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

type mappedWatchStream struct {
	events chan objectwatch.Event
	once   sync.Once
	done   chan struct{}
}

func newMappedWatchStream() *mappedWatchStream {
	return &mappedWatchStream{
		events: make(chan objectwatch.Event, 1),
		done:   make(chan struct{}),
	}
}

func (s *mappedWatchStream) send(t testing.TB, event objectwatch.Event) {
	t.Helper()

	select {
	case s.events <- event:
	case <-time.After(5 * time.Second):
		t.Fatal("watch stream did not accept event")
	}
}

func (s *mappedWatchStream) Next(ctx context.Context) (objectwatch.Event, error) {
	select {
	case event := <-s.events:
		return event, nil
	case <-ctx.Done():
		s.once.Do(func() { close(s.done) })
		return objectwatch.Event{}, ctx.Err()
	}
}

func (s *mappedWatchStream) Close() error {
	return nil
}

type mappedRecordingReconciler struct {
	mu      sync.Mutex
	changed chan struct{}
	records []mappedRecord
}

func newMappedRecordingReconciler() *mappedRecordingReconciler {
	return &mappedRecordingReconciler{changed: make(chan struct{})}
}

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

func (r *mappedRecordingReconciler) recordCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.records)
}

type mappedRecord struct {
	request  objectreconciler.Request
	snapshot objectreconciler.Snapshot
}

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

func mappedChangedEvent(t testing.TB, change objectstore.Change) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Changed(change)
	requireNoError(t, err)

	return event
}

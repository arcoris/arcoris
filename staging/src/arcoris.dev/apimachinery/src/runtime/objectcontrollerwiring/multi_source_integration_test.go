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
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectcache"
	"arcoris.dev/apimachinery/runtime/objectenqueue"
	"arcoris.dev/apimachinery/runtime/objectindex"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestMultiSourcePrimarySameObjectInputReconcilesFromPrimarySnapshot(t *testing.T) {
	targetKey := runTestKey("target")
	reconciler := newMappedRecordingReconciler()
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newMultiSourceIntegrationGraph(t, primary, nil, objectenqueue.ListedObject(), objectenqueue.ChangedObject(), reconciler)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := runMultiSourceAsync(ctx, graph)
	waitForSignal(t, primary.watchCalled)

	records := reconciler.waitForRecords(t, 1)
	requireMappedRecordKeys(t, records, targetKey)
	requireMappedRecordSnapshotContains(t, records[0], targetKey, 1)

	cancel()
	requireMultiSourceRunResult(t, result, context.Canceled)
}

func TestMultiSourceSecondaryMappedInputFeedsSharedController(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	reconciler := newMappedRecordingReconciler()
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	secondary := &runTestListerWatcher{
		read:        runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{mappedInputConfig(secondary, listItemMapperForKeys(targetKey), zeroChangeMapper())},
		zeroListItemMapper(),
		zeroChangeMapper(),
		reconciler,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := runMultiSourceAsync(ctx, graph)
	waitForSignal(t, primary.watchCalled)
	waitForSignal(t, secondary.watchCalled)

	records := reconciler.waitForRecords(t, 1)
	requireMappedRecordKeys(t, records, targetKey)
	requireMappedRecordSnapshotContains(t, records[0], targetKey, 1)
	requireMappedRecordSnapshotMissing(t, records[0], sourceKey)
	requireCacheContains(t, graph.Secondary()[0].Cache(), sourceKey, 2)

	cancel()
	requireMultiSourceRunResult(t, result, context.Canceled)
}

func TestMultiSourceSecondaryMapperUsesInputIndex(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	secondaryIndex := newDesiredObjectIndex(t)
	reconciler := newMappedRecordingReconciler()
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	secondary := &runTestListerWatcher{
		read:        runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	secondaryConfig := mappedInputConfig(
		secondary,
		listItemMapperUsingIndexToEmitTarget(secondaryIndex, "desired", "source", targetKey),
		zeroChangeMapper(),
	)
	secondaryConfig.Indexes = []*objectindex.Index{secondaryIndex}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{secondaryConfig},
		zeroListItemMapper(),
		zeroChangeMapper(),
		reconciler,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := runMultiSourceAsync(ctx, graph)
	waitForSignal(t, primary.watchCalled)
	waitForSignal(t, secondary.watchCalled)

	records := reconciler.waitForRecords(t, 1)
	requireMappedRecordKeys(t, records, targetKey)
	requireMappedRecordSnapshotContains(t, records[0], targetKey, 1)
	requireMappedRecordSnapshotMissing(t, records[0], sourceKey)
	requireCacheContains(t, graph.Secondary()[0].Cache(), sourceKey, 2)
	requireObjectIndexKeys(t, secondaryIndex, "desired", "source", sourceKey)

	cancel()
	requireMultiSourceRunResult(t, result, context.Canceled)
}

func TestMultiSourceSecondaryChangedInputMapsToPrimarySnapshot(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	stream := newMappedWatchStream()
	reconciler := newMappedRecordingReconciler()
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	secondary := &runTestListerWatcher{
		read:        runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
		stream:      stream,
		watchCalled: make(chan struct{}),
	}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{mappedInputConfig(secondary, zeroListItemMapper(), changeMapperForKeys(targetKey))},
		zeroListItemMapper(),
		zeroChangeMapper(),
		reconciler,
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result := runMultiSourceAsync(ctx, graph)
	waitForSignal(t, primary.watchCalled)
	waitForSignal(t, secondary.watchCalled)

	stream.send(t, mappedChangedEvent(t, updatedMappedChange(t, sourceKey, 2, 3)))
	records := reconciler.waitForRecords(t, 1)
	requireMappedRecordKeys(t, records, targetKey)
	requireMappedRecordSnapshotContains(t, records[0], targetKey, 1)
	requireCacheContains(t, graph.Secondary()[0].Cache(), sourceKey, 3)

	cancel()
	requireMultiSourceRunResult(t, result, context.Canceled)
}

func TestMultiSourceSharedQueueDeduplicatesPendingWork(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	watchErr := errors.New("watch opened")
	primary := &runTestListerWatcher{
		read:     runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		watchErr: watchErr,
	}
	secondary := &runTestListerWatcher{
		read:     runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
		watchErr: watchErr,
	}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{mappedInputConfig(secondary, listItemMapperForKeys(targetKey), zeroChangeMapper())},
		objectenqueue.ListedObject(),
		objectenqueue.ChangedObject(),
		newMappedRecordingReconciler(),
	)

	requireErrorIs(t, graph.Primary().Reflector().Run(context.Background()), watchErr)
	requireErrorIs(t, graph.Secondary()[0].Reflector().Run(context.Background()), watchErr)

	if got := graph.Queue().Len(); got != 1 {
		t.Fatalf("queue length = %d; want 1", got)
	}
	graph.Queue().ShutDown()
	item, err := graph.Queue().Get(context.Background())
	requireNoError(t, err)
	if !item.Key.Equal(targetKey) {
		t.Fatalf("queued key = %#v; want %#v", item.Key, targetKey)
	}
	requireNoError(t, graph.Queue().Done(item))
	_, err = graph.Queue().Get(context.Background())
	requireErrorIs(t, err, objectworkqueue.ErrShutDown)
}

func TestMultiSourceAllowsZeroMappedWork(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	reconciler := newMappedRecordingReconciler()
	primary := &runTestListerWatcher{
		read:        runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	secondary := &runTestListerWatcher{
		read:        runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
		stream:      runTestWaitingStream(),
		watchCalled: make(chan struct{}),
	}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{mappedInputConfig(secondary, zeroListItemMapper(), zeroChangeMapper())},
		zeroListItemMapper(),
		zeroChangeMapper(),
		reconciler,
	)
	ctx, cancel := context.WithCancel(context.Background())

	result := runMultiSourceAsync(ctx, graph)
	waitForSignal(t, primary.watchCalled)
	waitForSignal(t, secondary.watchCalled)
	cancel()

	requireMultiSourceRunResult(t, result, context.Canceled)
	if count := reconciler.recordCount(); count != 0 {
		t.Fatalf("reconciler records = %d; want 0", count)
	}
}

func TestMultiSourceReturnsSecondaryListedMapperError(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	mapperErr := errors.New("secondary listed mapper failed")
	primaryStream := runTestWaitingStream()
	primary := &runTestListerWatcher{
		read:   runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream: primaryStream,
	}
	secondary := &runTestListerWatcher{
		read: runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
	}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{mappedInputConfig(secondary, listItemMapperError(mapperErr), zeroChangeMapper())},
		zeroListItemMapper(),
		zeroChangeMapper(),
		newMappedRecordingReconciler(),
	)

	err := RunMultiSource(context.Background(), graph)

	requireErrorIs(t, err, mapperErr)
	if !graph.Queue().IsShutDown() {
		t.Fatal("queue is not shut down")
	}
	waitForSignal(t, primaryStream.done)
}

func TestMultiSourceReturnsSecondaryChangedMapperError(t *testing.T) {
	targetKey := runTestKey("target")
	sourceKey := runTestKey("source")
	mapperErr := errors.New("secondary changed mapper failed")
	primaryStream := runTestWaitingStream()
	stream := newMappedWatchStream()
	stream.send(t, mappedChangedEvent(t, updatedMappedChange(t, sourceKey, 2, 3)))
	primary := &runTestListerWatcher{
		read:   runTestRead(t, 1, runTestItem(targetKey, 1, "target")),
		stream: primaryStream,
	}
	secondary := &runTestListerWatcher{
		read:   runTestRead(t, 2, runTestItem(sourceKey, 2, "source")),
		stream: stream,
	}
	graph := newMultiSourceIntegrationGraph(
		t,
		primary,
		[]InputConfig{mappedInputConfig(secondary, zeroListItemMapper(), changeMapperError(mapperErr))},
		zeroListItemMapper(),
		zeroChangeMapper(),
		newMappedRecordingReconciler(),
	)

	err := RunMultiSource(context.Background(), graph)

	requireErrorIs(t, err, mapperErr)
	if !graph.Queue().IsShutDown() {
		t.Fatal("queue is not shut down")
	}
	waitForSignal(t, primaryStream.done)
}

func newMultiSourceIntegrationGraph(
	t testing.TB,
	primary *runTestListerWatcher,
	secondary []InputConfig,
	primaryListed objectenqueue.ListItemMapper,
	primaryChanged objectenqueue.Mapper,
	reconciler *mappedRecordingReconciler,
) *MultiSource {
	t.Helper()

	config := validMultiSourceConfig()
	config.Primary = mappedInputConfig(primary, primaryListed, primaryChanged)
	config.Secondary = secondary
	config.Reconciler = reconciler

	graph, err := NewMultiSource(config)
	requireNoError(t, err)

	return graph
}

func mappedInputConfig(
	source *runTestListerWatcher,
	listed objectenqueue.ListItemMapper,
	changed objectenqueue.Mapper,
) InputConfig {
	return InputConfig{
		Source:     source,
		Collection: runTestCollection(),
		Listed:     listed,
		Changed:    changed,
	}
}

func listItemMapperUsingIndexToEmitTarget(
	index *objectindex.Index,
	name objectindex.Name,
	value objectindex.Value,
	target objectstore.Key,
) objectenqueue.ListItemMapper {
	return objectenqueue.ListItemMapperFunc(func(_ objectstore.ListItem, emit objectenqueue.EmitFunc) error {
		keys, err := index.Lookup(name, value)
		if err != nil {
			return err
		}
		if len(keys) == 0 {
			return errors.New("index lookup did not find source key")
		}

		return emitKeys([]objectstore.Key{target}, emit)
	})
}

func requireMappedRecordSnapshotMissing(t testing.TB, record mappedRecord, key objectstore.Key) {
	t.Helper()

	result, err := record.snapshot.View.Get(key)
	requireNoError(t, err)
	if result.Found {
		t.Fatalf("snapshot revision %s contains %#v; want absent", record.snapshot.Revision, key)
	}
}

func requireCacheContains(
	t testing.TB,
	cache *objectcache.Cache,
	key objectstore.Key,
	revision objectstore.Revision,
) {
	t.Helper()

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	result, err := snap.Value.Get(key)
	requireNoError(t, err)
	if !result.Found {
		t.Fatalf("cache snapshot revision %s does not contain %#v", snap.Revision, key)
	}
	if result.State.Revision != revision {
		t.Fatalf("cache state revision = %s; want %s", result.State.Revision, revision)
	}
}

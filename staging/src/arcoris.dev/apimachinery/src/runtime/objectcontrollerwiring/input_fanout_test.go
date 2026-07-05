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

func TestNewInputFanoutRejectsNilIndex(t *testing.T) {
	cache, err := objectcache.New(runTestCollection())
	requireNoError(t, err)
	enqueueSink := newNoopReflectorSink(t)

	_, _, err = newInputFanout(cache, []*objectindex.Index{nil}, enqueueSink)

	requireErrorIs(t, err, ErrNilIndex)
}

func TestNewInputFanoutUpdatesCacheAndIndexesBeforeEnqueueOnReplace(t *testing.T) {
	key := runTestKey("source")
	cache, err := objectcache.New(runTestCollection())
	requireNoError(t, err)
	desiredIndex := newDesiredObjectIndex(t)
	nameIndex := newNameObjectIndex(t)
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	enqueueSink, err := objectenqueue.NewReflectorSink(objectenqueue.ReflectorSinkConfig{
		Queue: queue,
		Listed: objectenqueue.ListItemMapperFunc(
			func(_ objectstore.ListItem, emit objectenqueue.EmitFunc) error {
				requireCacheContains(t, cache, key, 1)
				requireObjectIndexKeys(t, desiredIndex, "desired", "source", key)
				requireObjectIndexKeys(t, nameIndex, "name", "source", key)

				return emit(objectworkqueue.Item{Key: key})
			},
		),
		Changed: zeroChangeMapper(),
	})
	requireNoError(t, err)
	fanout, indexes, err := newInputFanout(
		cache,
		[]*objectindex.Index{desiredIndex, nameIndex},
		enqueueSink,
	)
	requireNoError(t, err)
	if len(indexes) != 2 || indexes[0] != desiredIndex || indexes[1] != nameIndex {
		t.Fatalf("indexes = %#v; want configured index order", indexes)
	}

	requireNoError(t, fanout.Replace(context.Background(), runTestRead(t, 1, runTestItem(key, 1, "source"))))

	requireQueueItem(t, queue, key)
}

func TestNewInputFanoutUpdatesIndexBeforeEnqueueOnApplyChange(t *testing.T) {
	key := runTestKey("source")
	cache, err := objectcache.New(runTestCollection())
	requireNoError(t, err)
	index := newDesiredObjectIndex(t)
	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 4})
	requireNoError(t, err)
	enqueueSink, err := objectenqueue.NewReflectorSink(objectenqueue.ReflectorSinkConfig{
		Queue:  queue,
		Listed: zeroListItemMapper(),
		Changed: objectenqueue.MapperFunc(
			func(_ objectstore.Change, emit objectenqueue.EmitFunc) error {
				requireObjectIndexKeys(t, index, "desired", "source")
				requireObjectIndexKeys(t, index, "desired", "source-updated", key)
				requireCacheContains(t, cache, key, 2)

				return emit(objectworkqueue.Item{Key: key})
			},
		),
	})
	requireNoError(t, err)
	fanout, _, err := newInputFanout(cache, []*objectindex.Index{index}, enqueueSink)
	requireNoError(t, err)
	requireNoError(t, fanout.Replace(context.Background(), runTestRead(t, 1, runTestItem(key, 1, "source"))))

	requireNoError(t, fanout.ApplyChange(context.Background(), updatedMappedChange(t, key, 1, 2)))

	requireQueueItem(t, queue, key)
}

func TestCopyIndexesReturnsDetachedSlice(t *testing.T) {
	index := newDesiredObjectIndex(t)
	original := []*objectindex.Index{index}
	indexes, err := copyIndexes(original)
	requireNoError(t, err)

	original[0] = nil
	if indexes[0] != index {
		t.Fatal("copyIndexes reused caller slice")
	}

	indexes[0] = nil

	again, err := copyIndexes([]*objectindex.Index{index})
	requireNoError(t, err)
	if again[0] != index {
		t.Fatal("copyIndexes did not preserve configured index pointer")
	}
}

func newNoopReflectorSink(t testing.TB) *objectenqueue.ReflectorSink {
	t.Helper()

	queue, err := objectworkqueue.New(objectworkqueue.Options{Capacity: 1})
	requireNoError(t, err)
	sink, err := objectenqueue.NewReflectorSink(objectenqueue.ReflectorSinkConfig{
		Queue:   queue,
		Listed:  zeroListItemMapper(),
		Changed: zeroChangeMapper(),
	})
	requireNoError(t, err)

	return sink
}

func newDesiredObjectIndex(t testing.TB) *objectindex.Index {
	t.Helper()

	index, err := objectindex.New(objectindex.Definition{
		Name:      "desired",
		Extractor: desiredObjectExtractor(),
	})
	requireNoError(t, err)

	return index
}

func newNameObjectIndex(t testing.TB) *objectindex.Index {
	t.Helper()

	index, err := objectindex.New(objectindex.Definition{
		Name: "name",
		Extractor: objectindex.ExtractorFunc(
			func(item objectstore.ListItem, emit objectindex.EmitFunc) error {
				return emit(objectindex.Value(item.Key.Object.Name))
			},
		),
	})
	requireNoError(t, err)

	return index
}

func desiredObjectExtractor() objectindex.Extractor {
	return objectindex.ExtractorFunc(func(item objectstore.ListItem, emit objectindex.EmitFunc) error {
		desired, ok := item.State.Object.Desired.AsString()
		if !ok {
			return errors.New("desired is not string")
		}

		return emit(objectindex.Value(desired))
	})
}

func newFailingDesiredObjectIndex(t testing.TB, failValue string) (*objectindex.Index, error) {
	t.Helper()

	extractorErr := errors.New("index extractor failed")
	index, err := objectindex.New(objectindex.Definition{
		Name: "desired",
		Extractor: objectindex.ExtractorFunc(
			func(item objectstore.ListItem, emit objectindex.EmitFunc) error {
				desired, ok := item.State.Object.Desired.AsString()
				if !ok {
					return errors.New("desired is not string")
				}
				if desired == failValue {
					return extractorErr
				}

				return emit(objectindex.Value(desired))
			},
		),
	})
	requireNoError(t, err)

	return index, extractorErr
}

func requireObjectIndexKeys(
	t testing.TB,
	index *objectindex.Index,
	name objectindex.Name,
	value objectindex.Value,
	want ...objectstore.Key,
) {
	t.Helper()

	keys, err := index.Lookup(name, value)
	requireNoError(t, err)
	requireObjectKeys(t, keys, want...)
}

func requireObjectKeys(t testing.TB, got []objectstore.Key, want ...objectstore.Key) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("keys = %d; want %d: %#v", len(got), len(want), got)
	}
	for i, key := range want {
		if !got[i].Equal(key) {
			t.Fatalf("key %d = %#v; want %#v", i, got[i], key)
		}
	}
}

func requireQueueItem(t testing.TB, queue *objectworkqueue.Queue, want objectstore.Key) {
	t.Helper()

	item, err := queue.Get(context.Background())
	requireNoError(t, err)
	if !item.Key.Equal(want) {
		t.Fatalf("queued key = %#v; want %#v", item.Key, want)
	}
	requireNoError(t, queue.Done(item))
}

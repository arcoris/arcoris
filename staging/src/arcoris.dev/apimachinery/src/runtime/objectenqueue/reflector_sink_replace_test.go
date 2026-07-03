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

package objectenqueue

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestReflectorSinkReplaceRejectsInvalidRead(t *testing.T) {
	sink := newTestReflectorSink(t, newSinkQueue())

	err := sink.Replace(context.Background(), storewatchapi.CollectionRead{})

	requireErrorIs(t, err, storewatchapi.ErrInvalidCollectionRead)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkReplaceEnqueuesCurrentMatchesInReadOrder(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")
	b := testListItem(2, 1, "b")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a, b)))

	queue.requireItems(t, objectworkqueue.Item{Key: a.Key}, objectworkqueue.Item{Key: b.Key})
	requireKnownKeys(t, sink, a.Key, b.Key)
}

func TestReflectorSinkReplaceZeroPredicateMatchesAll(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue, withPredicate(objectquery.Predicate{}))
	a := testListItem(1, 1, "a")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a)))

	queue.requireItems(t, objectworkqueue.Item{Key: a.Key})
	requireKnownKeys(t, sink, a.Key)
}

func TestReflectorSinkReplaceSkipsCurrentPredicateMiss(t *testing.T) {
	queue := newSinkQueue()
	predicate := mustPredicate(t, mustLabelEquals(t, "env", "match"))
	sink := newTestReflectorSink(t, queue, withPredicate(predicate))
	a := testListItem(1, 1, "miss")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a)))

	queue.requireItems(t)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkReplaceEnqueuesMissingAfterCurrentMatches(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")
	b1 := testListItem(2, 1, "b1")
	b2 := testListItem(2, 2, "b2")
	c := testListItem(3, 2, "c")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a, b1)))
	queue.reset()

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 2, b2, c)))

	queue.requireItems(t,
		objectworkqueue.Item{Key: b2.Key},
		objectworkqueue.Item{Key: c.Key},
		objectworkqueue.Item{Key: a.Key},
	)
	requireKnownKeys(t, sink, b2.Key, c.Key)
	requireKnownRevision(t, sink, b2.Key, 2)
}

func TestReflectorSinkReplaceEnqueuesQueryLeftUsingPreviousKnownItem(t *testing.T) {
	queue := newSinkQueue()
	predicate := mustPredicate(t, mustLabelEquals(t, "env", "keep"))
	var mapped []objectstore.ListItem
	listed := ListItemMapperFunc(func(item objectstore.ListItem, emit EmitFunc) error {
		mapped = append(mapped, item)
		return emit(objectworkqueue.Item{Key: item.Key})
	})
	sink := newTestReflectorSink(t, queue, withPredicate(predicate), withListed(listed))
	previous := testListItem(1, 1, "keep")
	current := testListItem(1, 2, "drop")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, previous)))
	queue.reset()
	mapped = nil

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 2, current)))

	queue.requireItems(t, objectworkqueue.Item{Key: previous.Key})
	if len(mapped) != 1 {
		t.Fatalf("mapped items = %d; want 1", len(mapped))
	}
	requireListItem(t, mapped[0], previous)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkReplacePreservesPreviousKnownOrderForMissing(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")
	b := testListItem(2, 1, "b")
	c := testListItem(3, 1, "c")

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, c, a, b)))
	queue.reset()

	requireNoError(t, sink.Replace(context.Background(), testRead(t, 2)))

	queue.requireItems(t,
		objectworkqueue.Item{Key: c.Key},
		objectworkqueue.Item{Key: a.Key},
		objectworkqueue.Item{Key: b.Key},
	)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkReplaceReturnsMapperErrorUnchanged(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue, withListed(ListItemMapperFunc(func(objectstore.ListItem, EmitFunc) error {
		return errMapperFailed
	})))
	a := testListItem(1, 1, "a")

	err := sink.Replace(context.Background(), testRead(t, 1, a))

	requireErrorSame(t, err, errMapperFailed)
	queue.requireItems(t)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkReplaceReturnsQueueErrorUnchanged(t *testing.T) {
	queue := newSinkQueue()
	queue.err = errQueueFailed
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")

	err := sink.Replace(context.Background(), testRead(t, 1, a))

	requireErrorSame(t, err, errQueueFailed)
	queue.requireItems(t, objectworkqueue.Item{Key: a.Key})
	requireKnownKeys(t, sink)
}

func TestReflectorSinkReplaceDoesNotMutateKnownOnLaterQueueError(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)
	a := testListItem(1, 1, "a")
	b := testListItem(2, 2, "b")
	c := testListItem(3, 2, "c")
	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, a)))
	queue.reset()
	queue.err = errQueueFailed
	queue.failAt = 2

	err := sink.Replace(context.Background(), testRead(t, 2, b, c))

	requireErrorSame(t, err, errQueueFailed)
	requireKnownKeys(t, sink, a.Key)
	requireKnownRevision(t, sink, a.Key, 1)
}

func TestReflectorSinkReplaceDoesNotCallChangedMapper(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue, withChanged(MapperFunc(func(objectstore.Change, EmitFunc) error {
		return errors.New("changed mapper should not be called")
	})))

	err := sink.Replace(context.Background(), testRead(t, 1, testListItem(1, 1, "a")))

	requireNoError(t, err)
}

func TestReflectorSinkReplacePassesContextUnchanged(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey{}, "ctx")
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)

	requireNoError(t, sink.Replace(ctx, testRead(t, 1, testListItem(1, 1, "a"))))

	queue.requireContexts(t, ctx)
}

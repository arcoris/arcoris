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
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestReflectorSinkApplyChangeReturnsProjectChangeError(t *testing.T) {
	sink := newTestReflectorSink(t, newSinkQueue())

	err := sink.ApplyChange(context.Background(), objectstore.Change{})

	requireErrorIs(t, err, objectquery.ErrInvalidChange)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkApplyChangeIgnoresIgnoredProjection(t *testing.T) {
	queue := newSinkQueue()
	predicate := mustPredicate(t, mustLabelEquals(t, "env", "match"))
	sink := newTestReflectorSink(t, queue, withPredicate(predicate))

	requireNoError(t, sink.ApplyChange(context.Background(), createdChange(t, 1)))

	queue.requireItems(t)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkApplyChangeEnqueuesEnteredAndAddsKnown(t *testing.T) {
	queue := newSinkQueue()
	predicate := mustPredicate(t, mustLabelEquals(t, "env", "created"))
	sink := newTestReflectorSink(t, queue, withPredicate(predicate))
	change := createdChange(t, 1)

	requireNoError(t, sink.ApplyChange(context.Background(), change))

	queue.requireItems(t, objectworkqueue.Item{Key: change.Key})
	requireKnownKeys(t, sink, change.Key)
	requireKnownRevision(t, sink, change.Key, change.After.Revision)
}

func TestReflectorSinkApplyChangeEnqueuesUpdatedAndUpdatesKnown(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)
	before := testListItem(1, 1, "before")
	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, before)))
	queue.reset()
	change := updatedChange(t, 1)

	requireNoError(t, sink.ApplyChange(context.Background(), change))

	queue.requireItems(t, objectworkqueue.Item{Key: change.Key})
	requireKnownKeys(t, sink, change.Key)
	requireKnownRevision(t, sink, change.Key, change.After.Revision)
}

func TestReflectorSinkApplyChangeEnqueuesLeftAndRemovesKnown(t *testing.T) {
	queue := newSinkQueue()
	predicate := mustPredicate(t, mustLabelEquals(t, "env", "before"))
	sink := newTestReflectorSink(t, queue, withPredicate(predicate))
	before := testListItem(1, 1, "before")
	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, before)))
	queue.reset()
	change := updatedChange(t, 1)

	requireNoError(t, sink.ApplyChange(context.Background(), change))

	queue.requireItems(t, objectworkqueue.Item{Key: change.Key})
	requireUnknown(t, sink, change.Key)
	requireKnownKeys(t, sink)
}

func TestReflectorSinkApplyChangeDoesNotMutateKnownOnMapperError(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue, withChanged(MapperFunc(func(objectstore.Change, EmitFunc) error {
		return errMapperFailed
	})))
	before := testListItem(1, 1, "before")
	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, before)))
	queue.reset()
	change := updatedChange(t, 1)

	err := sink.ApplyChange(context.Background(), change)

	requireErrorSame(t, err, errMapperFailed)
	queue.requireItems(t)
	requireKnownKeys(t, sink, before.Key)
	requireKnownRevision(t, sink, before.Key, before.State.Revision)
}

func TestReflectorSinkApplyChangeDoesNotMutateKnownOnQueueError(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)
	before := testListItem(1, 1, "before")
	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, before)))
	queue.reset()
	queue.err = errQueueFailed
	change := updatedChange(t, 1)

	err := sink.ApplyChange(context.Background(), change)

	requireErrorSame(t, err, errQueueFailed)
	queue.requireItems(t, objectworkqueue.Item{Key: change.Key})
	requireKnownKeys(t, sink, before.Key)
	requireKnownRevision(t, sink, before.Key, before.State.Revision)
}

func TestReflectorSinkApplyChangePreservesKnownOnFailedLeft(t *testing.T) {
	queue := newSinkQueue()
	predicate := mustPredicate(t, mustLabelEquals(t, "env", "before"))
	sink := newTestReflectorSink(t, queue, withPredicate(predicate))
	before := testListItem(1, 1, "before")
	requireNoError(t, sink.Replace(context.Background(), testRead(t, 1, before)))
	queue.reset()
	queue.err = errQueueFailed
	change := updatedChange(t, 1)

	err := sink.ApplyChange(context.Background(), change)

	requireErrorSame(t, err, errQueueFailed)
	requireKnownKeys(t, sink, before.Key)
}

func TestReflectorSinkApplyChangeDoesNotCallListedMapper(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue, withListed(ListItemMapperFunc(func(objectstore.ListItem, EmitFunc) error {
		return errors.New("listed mapper should not be called")
	})))

	err := sink.ApplyChange(context.Background(), createdChange(t, 1))

	requireNoError(t, err)
	queue.requireItems(t, objectworkqueue.Item{Key: testKey(1)})
}

func TestReflectorSinkApplyChangePassesContextUnchanged(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey{}, "ctx")
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)

	requireNoError(t, sink.ApplyChange(ctx, createdChange(t, 1)))

	queue.requireContexts(t, ctx)
}

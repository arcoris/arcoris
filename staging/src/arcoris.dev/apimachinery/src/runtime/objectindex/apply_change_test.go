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

package objectindex

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestApplyChangeAddsCreatedObjectMemberships(t *testing.T) {
	key := testKey("task-1")
	index := newDesiredIndex(t)

	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))

	requireLookupKeys(t, index, "desired", "worker-a", key)
}

func TestApplyChangeUpdatesMembershipWhenValueChanges(t *testing.T) {
	key := testKey("task-1")
	index := newDesiredIndex(t)
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))

	requireNoError(t, index.ApplyChange(context.Background(), updatedChange(key, 1, "worker-a", 2, "worker-b")))

	requireLookupKeys(t, index, "desired", "worker-a")
	requireLookupKeys(t, index, "desired", "worker-b", key)
}

func TestApplyChangePreservesUnchangedMembershipOrder(t *testing.T) {
	task1 := testKey("task-1")
	task2 := testKey("task-2")
	index := newDesiredIndex(t)
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(task1, 1, "worker-a")))
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(task2, 2, "worker-a")))

	requireNoError(t, index.ApplyChange(context.Background(), updatedChange(task1, 1, "worker-a", 3, "worker-a")))

	requireLookupKeys(t, index, "desired", "worker-a", task1, task2)
}

func TestApplyChangeRemovesMembershipWhenObjectNoLongerEmitsValue(t *testing.T) {
	key := testKey("task-1")
	index, err := New(Definition{
		Name: "desired",
		Extractor: ExtractorFunc(func(item objectstore.ListItem, emit EmitFunc) error {
			desired, ok := item.State.Object.Desired.AsString()
			if !ok || desired == "none" {
				return nil
			}

			return emit(Value(desired))
		}),
	})
	requireNoError(t, err)
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))

	requireNoError(t, index.ApplyChange(context.Background(), updatedChange(key, 1, "worker-a", 2, "none")))

	requireLookupKeys(t, index, "desired", "worker-a")
}

func TestApplyChangeRemovesMembershipsOnDelete(t *testing.T) {
	key := testKey("task-1")
	index := newDesiredIndex(t)
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))

	requireNoError(t, index.ApplyChange(context.Background(), deletedChange(key, 1, "worker-a", 2)))

	requireLookupKeys(t, index, "desired", "worker-a")
}

func TestApplyChangeSupportsMultipleDefinitions(t *testing.T) {
	key := testKey("task-1")
	index, err := New(
		Definition{Name: "desired", Extractor: desiredExtractor()},
		Definition{Name: "name", Extractor: nameExtractor()},
	)
	requireNoError(t, err)

	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))

	requireLookupKeys(t, index, "desired", "worker-a", key)
	requireLookupKeys(t, index, "name", "task-1", key)
}

func TestApplyChangeIsAllOrNothingWhenExtractorFails(t *testing.T) {
	key := testKey("task-1")
	extractorErr := errors.New("extract failed")
	index := newDesiredIndex(t)
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))
	index.definitions["desired"] = errorExtractor(extractorErr)

	err := index.ApplyChange(context.Background(), updatedChange(key, 1, "worker-a", 2, "worker-b"))

	requireErrorIs(t, err, extractorErr)
	requireLookupKeys(t, index, "desired", "worker-a", key)
	requireLookupKeys(t, index, "desired", "worker-b")
}

func TestApplyChangeRespectsContextCancellation(t *testing.T) {
	key := testKey("task-1")
	index := newDesiredIndex(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := index.ApplyChange(ctx, createdChange(key, 1, "worker-a"))

	requireErrorIs(t, err, context.Canceled)
	requireLookupKeys(t, index, "desired", "worker-a")
}

func TestApplyChangeRejectsInvalidChange(t *testing.T) {
	index := newDesiredIndex(t)

	err := index.ApplyChange(context.Background(), objectstore.Change{})

	requireErrorIs(t, err, objectstore.ErrInvalidChange)
}

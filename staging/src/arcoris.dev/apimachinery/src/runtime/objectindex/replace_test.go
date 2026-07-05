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
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

func TestReplaceBuildsIndexFromListItems(t *testing.T) {
	task1 := testKey("task-1")
	task2 := testKey("task-2")
	index := newDesiredIndex(t)

	requireNoError(t, index.Replace(context.Background(), testRead(t, 2,
		testItem(task1, 1, "worker-a"),
		testItem(task2, 2, "worker-a"),
	)))

	requireLookupKeys(t, index, "desired", "worker-a", task1, task2)
}

func TestReplaceSupportsMultipleDefinitions(t *testing.T) {
	key := testKey("task-1")
	index, err := New(
		Definition{Name: "desired", Extractor: desiredExtractor()},
		Definition{Name: "name", Extractor: nameExtractor()},
	)
	requireNoError(t, err)

	requireNoError(t, index.Replace(context.Background(), testRead(t, 1, testItem(key, 1, "worker-a"))))

	requireLookupKeys(t, index, "desired", "worker-a", key)
	requireLookupKeys(t, index, "name", "task-1", key)
}

func TestReplaceSupportsMultipleAndZeroValues(t *testing.T) {
	multi := testKey("task-multi")
	zero := testKey("task-zero")
	index, err := New(Definition{
		Name: "worker",
		Extractor: ExtractorFunc(func(item objectstore.ListItem, emit EmitFunc) error {
			if item.Key.Equal(zero) {
				return nil
			}
			if err := emit("worker-a"); err != nil {
				return err
			}

			return emit("worker-b")
		}),
	})
	requireNoError(t, err)

	requireNoError(t, index.Replace(context.Background(), testRead(t, 2,
		testItem(multi, 1, "ignored"),
		testItem(zero, 2, "ignored"),
	)))

	requireLookupKeys(t, index, "worker", "worker-a", multi)
	requireLookupKeys(t, index, "worker", "worker-b", multi)
	requireLookupKeys(t, index, "worker", "missing")
}

func TestReplaceCollapsesDuplicateValuesForOneObject(t *testing.T) {
	key := testKey("task-1")
	index, err := New(Definition{Name: "worker", Extractor: fixedExtractor("worker-a", "worker-a")})
	requireNoError(t, err)

	requireNoError(t, index.Replace(context.Background(), testRead(t, 1, testItem(key, 1, "ignored"))))

	requireLookupKeys(t, index, "worker", "worker-a", key)
}

func TestReplacePreservesListOrder(t *testing.T) {
	task1 := testKey("task-1")
	task2 := testKey("task-2")
	task3 := testKey("task-3")
	index, err := New(Definition{Name: "worker", Extractor: fixedExtractor("worker-a")})
	requireNoError(t, err)

	requireNoError(t, index.Replace(context.Background(), testRead(t, 3,
		testItem(task2, 2, "ignored"),
		testItem(task1, 1, "ignored"),
		testItem(task3, 3, "ignored"),
	)))

	requireLookupKeys(t, index, "worker", "worker-a", task2, task1, task3)
}

func TestReplaceIsAllOrNothingWhenExtractorFails(t *testing.T) {
	task1 := testKey("task-1")
	task2 := testKey("task-2")
	extractorErr := errors.New("extract failed")
	index := newDesiredIndex(t)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 1, testItem(task1, 1, "worker-a"))))
	index.definitions["desired"] = errorExtractor(extractorErr)

	err := index.Replace(context.Background(), testRead(t, 2, testItem(task2, 2, "worker-b")))

	requireErrorIs(t, err, extractorErr)
	requireLookupKeys(t, index, "desired", "worker-a", task1)
	requireLookupKeys(t, index, "desired", "worker-b")
}

func TestReplaceRespectsContextCancellation(t *testing.T) {
	index := newDesiredIndex(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := index.Replace(ctx, testRead(t, 1, testItem(testKey("task-1"), 1, "worker-a")))

	requireErrorIs(t, err, context.Canceled)
	requireLookupKeys(t, index, "desired", "worker-a")
}

func TestReplaceRejectsInvalidRead(t *testing.T) {
	index := newDesiredIndex(t)

	err := index.Replace(context.Background(), storewatchapi.CollectionRead{})

	requireErrorIs(t, err, storewatchapi.ErrInvalidCollectionRead)
}

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
)

func TestLookupReturnsMatchingKeys(t *testing.T) {
	task1 := testKey("task-1")
	task2 := testKey("task-2")
	index := newDesiredIndex(t)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 2,
		testItem(task1, 1, "worker-a"),
		testItem(task2, 2, "worker-a"),
	)))

	requireLookupKeys(t, index, "desired", "worker-a", task1, task2)
}

func TestLookupReturnsDetachedSlice(t *testing.T) {
	task1 := testKey("task-1")
	task2 := testKey("task-2")
	index := newDesiredIndex(t)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 2,
		testItem(task1, 1, "worker-a"),
		testItem(task2, 2, "worker-a"),
	)))

	keys, err := index.Lookup("desired", "worker-a")
	requireNoError(t, err)
	keys[0] = testKey("mutated")

	requireLookupKeys(t, index, "desired", "worker-a", task1, task2)
}

func TestLookupReturnsEmptyForKnownIndexValueWithoutKeys(t *testing.T) {
	index := newDesiredIndex(t)

	keys, err := index.Lookup("desired", "missing")

	requireNoError(t, err)
	if len(keys) != 0 {
		t.Fatalf("keys = %#v; want empty", keys)
	}
}

func TestLookupReturnsErrUnknownIndex(t *testing.T) {
	index := newDesiredIndex(t)

	_, err := index.Lookup("missing", "worker-a")

	requireErrorIs(t, err, ErrUnknownIndex)
}

func TestLookupRejectsNilIndex(t *testing.T) {
	var index *Index

	_, err := index.Lookup("desired", "worker-a")

	requireErrorIs(t, err, ErrInvalidIndex)
}

func TestLookupRejectsEmptyNameOrValue(t *testing.T) {
	index := newDesiredIndex(t)

	_, err := index.Lookup("", "worker-a")
	requireErrorIs(t, err, ErrInvalidIndex)

	_, err = index.Lookup("desired", "")
	requireErrorIs(t, err, ErrInvalidIndex)
}

func TestLookupIsSafeAfterFailedReplace(t *testing.T) {
	key := testKey("task-1")
	index := newDesiredIndex(t)
	requireNoError(t, index.Replace(context.Background(), testRead(t, 1, testItem(key, 1, "worker-a"))))
	index.definitions["desired"] = errorExtractor(errors.New("extract failed"))

	_ = index.Replace(context.Background(), testRead(t, 2, testItem(testKey("task-2"), 2, "worker-b")))

	requireLookupKeys(t, index, "desired", "worker-a", key)
}

func TestLookupIsSafeAfterFailedApplyChange(t *testing.T) {
	key := testKey("task-1")
	index := newDesiredIndex(t)
	requireNoError(t, index.ApplyChange(context.Background(), createdChange(key, 1, "worker-a")))
	index.definitions["desired"] = errorExtractor(errors.New("extract failed"))

	_ = index.ApplyChange(context.Background(), updatedChange(key, 1, "worker-a", 2, "worker-b"))

	requireLookupKeys(t, index, "desired", "worker-a", key)
}

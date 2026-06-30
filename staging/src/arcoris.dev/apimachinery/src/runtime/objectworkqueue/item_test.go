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

package objectworkqueue

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestItemCarriesObjectStoreKey(t *testing.T) {
	key := testKey(1)
	item := Item{Key: key}

	if !item.Key.Equal(key) {
		t.Fatalf("item key = %s; want %s", item.Key, key)
	}
}

func TestValidateItemRejectsZeroItem(t *testing.T) {
	err := validateItem(Item{})

	requireErrorIs(t, err, ErrInvalidItem)
	requireErrorIs(t, err, objectstore.ErrInvalidKey)
	if !errors.Is(err, ErrInvalidItem) {
		t.Fatalf("errors.Is(%v, ErrInvalidItem) = false; want true", err)
	}
}

func TestValidateItemAcceptsValidItem(t *testing.T) {
	requireNoError(t, validateItem(testItem(1)))
}

func TestAddRejectsZeroItemAndDoesNotConsumeCapacity(t *testing.T) {
	queue := newTestQueue(t, 1)

	requireErrorIs(t, queue.Add(context.Background(), Item{}), ErrInvalidItem)

	requireStats(t, queue, 0, 0)
	requireNoError(t, queue.Add(context.Background(), testItem(1)))
	requireStats(t, queue, 1, 0)
	requireInvariants(t, queue)
}

func TestTryAddRejectsZeroItemAndDoesNotConsumeCapacity(t *testing.T) {
	queue := newTestQueue(t, 1)

	requireErrorIs(t, queue.TryAdd(Item{}), ErrInvalidItem)

	requireStats(t, queue, 0, 0)
	requireNoError(t, queue.TryAdd(testItem(1)))
	requireStats(t, queue, 1, 0)
	requireInvariants(t, queue)
}

func TestDoneRejectsZeroItemAndDoesNotMutateState(t *testing.T) {
	queue := newTestQueue(t, 1)
	item := testItem(1)
	requireNoError(t, queue.Add(context.Background(), item))
	got, err := queue.Get(context.Background())
	requireNoError(t, err)
	requireItem(t, got, item)

	requireErrorIs(t, queue.Done(Item{}), ErrInvalidItem)

	requireStats(t, queue, 0, 1)
	requireNoError(t, queue.Done(item))
	requireStats(t, queue, 0, 0)
	requireInvariants(t, queue)
}

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

package objectstorewatch

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestStreamOverflowErrorPreservesContinuityAndOverflow(t *testing.T) {
	err := streamOverflowError()

	if !errors.Is(err, objectwatch.ErrContinuityLost) {
		t.Fatalf("errors.Is(%v, %v) = false", err, objectwatch.ErrContinuityLost)
	}
	if !errors.Is(err, ErrStreamOverflow) {
		t.Fatalf("errors.Is(%v, %v) = false", err, ErrStreamOverflow)
	}
}

func TestSlowWatcherLosesContinuityOnOverflow(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(1))
	stream := watchAfter(t, store, testCollection(), 0)

	createObject(t, store, testKey("system", 1), "one")
	createObject(t, store, testKey("system", 2), "two")

	_, err := stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, ErrStreamOverflow)

	_, err = stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
}

func TestSlowWatcherOverflowDoesNotPoisonHealthyWatcher(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(1))
	slow := watchAfter(t, store, testCollection(), 0)
	healthy := watchAfter(t, store, testCollection(), 0)

	first := createObject(t, store, testKey("system", 1), "one")
	if event := nextEvent(t, healthy); event.Revision != first.Revision {
		t.Fatalf("healthy first revision = %s; want %s", event.Revision, first.Revision)
	}
	second := createObject(t, store, testKey("system", 2), "two")

	if event := nextEvent(t, healthy); event.Revision != second.Revision {
		t.Fatalf("healthy second revision = %s; want %s", event.Revision, second.Revision)
	}
	_, err := slow.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	requireErrorIs(t, err, ErrStreamOverflow)
}

func TestFutureWatchReplaysHistoryAfterWatcherOverflow(t *testing.T) {
	store := testRuntimeStore(t, WithMaxHistory(10), WithStreamBuffer(1))
	slow := watchAfter(t, store, testCollection(), 0)

	first := createObject(t, store, testKey("system", 1), "one")
	second := createObject(t, store, testKey("system", 2), "two")
	_, err := slow.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrContinuityLost)

	replay := watchAfter(t, store, testCollection(), 0)
	for _, revision := range []objectstore.Revision{first.Revision, second.Revision} {
		if event := nextEvent(t, replay); event.Revision != revision {
			t.Fatalf("replay revision = %s; want %s", event.Revision, revision)
		}
	}
}

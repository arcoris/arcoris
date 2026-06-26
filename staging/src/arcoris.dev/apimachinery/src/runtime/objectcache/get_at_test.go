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

package objectcache

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestGetAtRejectsDisabledOutsideAndFutureReads(t *testing.T) {
	var nilCache *Cache
	key := testKey("system", 1)
	_, err := nilCache.GetAt(key, 1)
	requireErrorIs(t, err, ErrInvalidCache)

	cache := readyCache(t, 1, listItem(key, 1, "current"))

	_, err = cache.GetAt(key, 1)
	requireErrorIs(t, err, ErrHistoryDisabled)

	history := readyHistoryCache(t, 1, 1, listItem(key, 1, "current"))
	_, err = history.GetAt(otherResourceKey("system", 1), 1)
	requireErrorIs(t, err, ErrOutsideCollection)
	_, err = history.GetAt(key, 2)
	requireErrorIs(t, err, ErrFutureRevision)
}

func TestGetAtReturnsExactAndNearestRetainedLiveVersion(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 3, 1, listItem(key, 1, "one"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "one"), testState(key, 4, "four")),
	))
	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 4, "four"), testState(key, 9, "nine")),
	))

	exact, err := cache.GetAt(key, 4)
	requireNoError(t, err)
	if !exact.Found || exact.Revision != 4 || desiredString(t, exact.State) != "four" {
		t.Fatalf("GetAt(4) = %#v; want live four at revision 4", exact)
	}

	nearest, err := cache.GetAt(key, 8)
	requireNoError(t, err)
	if !nearest.Found || nearest.Revision != 8 || desiredString(t, nearest.State) != "four" {
		t.Fatalf("GetAt(8) = %#v; want live four served at revision 8", nearest)
	}
	if nearest.State.Revision != 4 {
		t.Fatalf("State.Revision = %s; want object revision 4", nearest.State.Revision)
	}
}

func TestGetAtReportsKnownAbsenceAfterRetainedTombstone(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "live"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "live"), 5),
	))

	result, err := cache.GetAt(key, 5)
	requireNoError(t, err)
	if result.Found || result.Revision != 5 {
		t.Fatalf("GetAt(delete) = %#v; want known absence at delete revision", result)
	}

	result, err = cache.GetAt(key, 6)
	requireErrorIs(t, err, ErrFutureRevision)
	if result.Found {
		t.Fatalf("future GetAt returned object: %#v", result)
	}
}

func TestGetAtReturnsPreviousLiveVersionBeforeDeleteWhenRetained(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "live"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "live"), 5),
	))

	result, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	if !result.Found || desiredString(t, result.State) != "live" {
		t.Fatalf("GetAt(before delete) = %#v; want retained live version", result)
	}
}

func TestGetAtCannotReturnLiveBeforeDeleteWhenOnlyTombstoneRetained(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 1, 1, listItem(key, 1, "live"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "live"), 5),
	))

	_, err := cache.GetAt(key, 1)
	requireErrorIs(t, err, ErrHistoryUnavailable)
}

func TestGetAtReturnsHistoryUnavailableWhenOlderThanPerObjectRetention(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 1, 1, listItem(key, 1, "one"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "one"), testState(key, 3, "three")),
	))

	_, err := cache.GetAt(key, 1)
	requireErrorIs(t, err, ErrHistoryUnavailable)
	current, err := cache.GetAt(key, 3)
	requireNoError(t, err)
	if !current.Found || desiredString(t, current.State) != "three" {
		t.Fatalf("GetAt(3) = %#v; want retained three", current)
	}
}

func TestGetAtAnswersCurrentKnownAbsenceForNeverSeenKey(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 4)

	result, err := cache.GetAt(key, 4)
	requireNoError(t, err)
	if result.Found || result.Revision != 4 {
		t.Fatalf("GetAt(current missing) = %#v; want known absence at revision 4", result)
	}

	_, err = cache.GetAt(key, 3)
	requireErrorIs(t, err, ErrHistoryUnavailable)
}

func TestGetAtReturnsDetachedHistoricalState(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 1, 1, listItem(key, 1, "stored"))

	result, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	mutateState(&result.State, "mutated")

	again, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	if desired := desiredString(t, again.State); desired != "stored" {
		t.Fatalf("desired = %q; want stored", desired)
	}
}

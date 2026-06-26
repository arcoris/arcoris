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

func TestPreviousLiveRejectsDisabledOutsideAndFutureReads(t *testing.T) {
	var nilCache *Cache
	key := testKey("system", 1)
	_, err := nilCache.PreviousLive(key, 1)
	requireErrorIs(t, err, ErrInvalidCache)

	cache := readyCache(t, 1, listItem(key, 1, "current"))

	_, err = cache.PreviousLive(key, 1)
	requireErrorIs(t, err, ErrHistoryDisabled)

	history := readyHistoryCache(t, 1, 1, listItem(key, 1, "current"))
	_, err = history.PreviousLive(otherResourceKey("system", 1), 1)
	requireErrorIs(t, err, ErrOutsideCollection)
	_, err = history.PreviousLive(key, 2)
	requireErrorIs(t, err, ErrFutureRevision)
}

func TestPreviousLiveReturnsNearestStrictlyBeforeRevision(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 3, 1, listItem(key, 1, "one"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "one"), testState(key, 5, "five")),
	))
	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 5, "five"), testState(key, 12, "twelve")),
	))

	result, err := cache.PreviousLive(key, 12)
	requireNoError(t, err)
	if !result.Found || result.Revision != 5 || desiredString(t, result.State) != "five" {
		t.Fatalf("PreviousLive(12) = %#v; want five at revision 5", result)
	}

	result, err = cache.PreviousLive(key, 11)
	requireNoError(t, err)
	if !result.Found || result.Revision != 5 || desiredString(t, result.State) != "five" {
		t.Fatalf("PreviousLive(11) = %#v; want five at revision 5", result)
	}
}

func TestPreviousLiveDoesNotComputeRevisionMinusOne(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 10, listItem(key, 3, "observed-at-ten"))
	other := testKey("system", 2)
	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustCreatedChange(other, testState(other, 11, "other")),
	))

	result, err := cache.PreviousLive(key, 11)
	requireNoError(t, err)
	if !result.Found || result.Revision != 10 || result.State.Revision != 3 {
		t.Fatalf("PreviousLive(11) = %#v; want replacement observation revision 10 and object revision 3", result)
	}
}

func TestPreviousLiveReturnsUnavailableWhenOnlyTombstoneIsRetained(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 1, 1, listItem(key, 1, "live"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "live"), 7),
	))

	_, err := cache.PreviousLive(key, 7)
	requireErrorIs(t, err, ErrHistoryUnavailable)
}

func TestPreviousLiveReturnsLiveVersionBeforeDeleteWhenRetained(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "live"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "live"), 7),
	))

	result, err := cache.PreviousLive(key, 7)
	requireNoError(t, err)
	if !result.Found || result.Revision != 1 || desiredString(t, result.State) != "live" {
		t.Fatalf("PreviousLive(delete) = %#v; want retained live version", result)
	}
}

func TestPreviousLiveReturnsDetachedState(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "live"))
	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "live"), testState(key, 3, "new")),
	))

	result, err := cache.PreviousLive(key, 3)
	requireNoError(t, err)
	mutateState(&result.State, "mutated")

	again, err := cache.PreviousLive(key, 3)
	requireNoError(t, err)
	if desired := desiredString(t, again.State); desired != "live" {
		t.Fatalf("desired = %q; want live", desired)
	}
}

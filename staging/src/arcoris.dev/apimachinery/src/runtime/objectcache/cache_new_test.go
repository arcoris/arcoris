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

func TestNewRejectsInvalidCollection(t *testing.T) {
	cache, err := New(objectstore.ListRequest{})

	if cache != nil {
		t.Fatalf("cache = %#v; want nil", cache)
	}
	requireErrorIs(t, err, ErrInvalidCache)
	requireErrorIs(t, err, objectstore.ErrInvalidListRequest)
}

func TestNewRejectsNilOption(t *testing.T) {
	cache, err := New(testCollection(), nil)

	if cache != nil {
		t.Fatalf("cache = %#v; want nil", cache)
	}
	requireErrorIs(t, err, ErrInvalidOption)
}

func TestNewStartsNotReady(t *testing.T) {
	cache, err := New(testCollection())
	requireNoError(t, err)

	if cache.Ready() {
		t.Fatalf("Ready() = true; want false")
	}
	_, err = cache.Revision()
	requireErrorIs(t, err, ErrNotReady)
	_, err = cache.List()
	requireErrorIs(t, err, ErrNotReady)
	_, err = cache.Get(testKey("system", 1))
	requireErrorIs(t, err, ErrNotReady)
	if got := cache.Len(); got != 0 {
		t.Fatalf("Len() = %d; want 0", got)
	}
}

func TestHistoricalReadsReturnNotReadyBeforeReplace(t *testing.T) {
	cache, err := New(testCollection(), WithHistory(HistoryPolicy{RetainedVersionsPerObject: 2}))
	requireNoError(t, err)
	key := testKey("system", 1)

	_, err = cache.GetAt(key, 1)
	requireErrorIs(t, err, ErrNotReady)
	_, err = cache.PreviousLive(key, 1)
	requireErrorIs(t, err, ErrNotReady)
}

func TestReplaceEmptyRevisionZeroMakesCacheReady(t *testing.T) {
	cache, err := New(testCollection())
	requireNoError(t, err)

	requireNoError(t, cache.Replace(context.Background(), collectionRead(t, testCollection(), 0)))

	if !cache.Ready() {
		t.Fatalf("Ready() = false; want true")
	}
	revision, err := cache.Revision()
	requireNoError(t, err)
	if revision != 0 {
		t.Fatalf("Revision() = %s; want 0", revision)
	}
	result, err := cache.List()
	requireNoError(t, err)
	if result.Revision != 0 || len(result.Items) != 0 {
		t.Fatalf("List() = %#v; want empty revision 0", result)
	}
}

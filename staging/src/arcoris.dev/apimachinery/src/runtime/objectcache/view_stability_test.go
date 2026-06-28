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

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestViewRemainsStableAfterCacheApplyChangeAndReplace(t *testing.T) {
	key := testKey("system", 1)
	replacement := testKey("system", 2)
	cache := readyCache(t, 1, listItem(key, 1, "initial"))

	snap, err := cache.ReadSnapshot()
	requireNoError(t, err)
	view := snap.Value

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "initial"), testState(key, 2, "updated")),
	))
	requireNoError(t, cache.Replace(
		context.Background(),
		collectionRead(t, testCollection(), 5, listItem(replacement, 5, "replacement")),
	))

	got, err := view.Get(key)
	requireNoError(t, err)
	if !got.Found || got.Revision != 1 || desiredString(t, got.State) != "initial" {
		t.Fatalf("view Get() = %#v; want initial object at revision 1", got)
	}

	list := view.List()
	if list.Revision != 1 {
		t.Fatalf("view List() revision = %s; want 1", list.Revision)
	}
	requireListOrder(t, list, key)

	initial := view.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "initial")
	}))
	requireListOrder(t, initial, key)

	updated := view.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "updated")
	}))
	if len(updated.Items) != 0 {
		t.Fatalf("view Query(updated) items = %d; want 0", len(updated.Items))
	}
}

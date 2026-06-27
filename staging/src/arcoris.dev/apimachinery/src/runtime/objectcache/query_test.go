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

	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestQueryRejectsNilOrNotReadyCache(t *testing.T) {
	var nilCache *Cache
	_, err := nilCache.Query(objectquery.Predicate{})
	requireErrorIs(t, err, ErrInvalidCache)

	cache, err := New(testCollection())
	requireNoError(t, err)
	_, err = cache.Query(objectquery.Predicate{})
	requireErrorIs(t, err, ErrNotReady)
}

func TestQueryZeroAndAllPredicatesMatchLatestItems(t *testing.T) {
	first := testKey("system", 2)
	second := testKey("system", 1)
	cache := readyCache(t, 9, listItem(first, 1, "first"), listItem(second, 2, "second"))

	zero, err := cache.Query(objectquery.Predicate{})
	requireNoError(t, err)
	if zero.Revision != 9 {
		t.Fatalf("zero predicate revision = %s; want 9", zero.Revision)
	}
	requireListOrder(t, zero, first, second)

	all, err := cache.Query(mustPredicate(t, objectquery.All()))
	requireNoError(t, err)
	requireListOrder(t, all, first, second)
}

func TestQueryNonePredicateMatchesNoItems(t *testing.T) {
	cache := readyCache(t, 3, listItem(testKey("system", 1), 1, "one"))

	result, err := cache.Query(mustPredicate(t, objectquery.None()))
	requireNoError(t, err)
	if result.Revision != 3 || len(result.Items) != 0 || result.Items != nil {
		t.Fatalf("Query(None) = %#v; want nil items at revision 3", result)
	}
}

func TestQueryFiltersByResourceKeyNamespaceAndObjectIdentity(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("other", 1)
	third := testKey("system", 2)
	cache := readyCache(
		t,
		6,
		listItem(first, 1, "first"),
		listItem(second, 2, "second"),
		listItem(third, 3, "third"),
	)

	resource, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.ResourceEquals(testResource())
	}))
	requireNoError(t, err)
	requireListOrder(t, resource, first, second, third)

	otherResourceResult, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.ResourceEquals(otherResource())
	}))
	requireNoError(t, err)
	if len(otherResourceResult.Items) != 0 {
		t.Fatalf("other resource items = %d; want 0", len(otherResourceResult.Items))
	}

	key, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.KeyEquals(third)
	}))
	requireNoError(t, err)
	requireListOrder(t, key, third)

	namespace, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.ObjectInNamespace("system")
	}))
	requireNoError(t, err)
	requireListOrder(t, namespace, first, third)

	name, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.ObjectWithName(metaidentity.Name("worker-1"))
	}))
	requireNoError(t, err)
	requireListOrder(t, name, first, second)

	object, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.ObjectEquals("system", metaidentity.Name("worker-1"))
	}))
	requireNoError(t, err)
	requireListOrder(t, object, first)
}

func TestQueryFiltersByLabelsAndAnnotations(t *testing.T) {
	first := testKey("system", 1)
	second := testKey("system", 2)
	cache := readyCache(
		t,
		4,
		listItemWithAnnotations(first, 1, "blue", map[string]string{"tier": "frontend"}),
		listItemWithAnnotations(second, 2, "green", map[string]string{"tier": "backend"}),
	)

	label, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "blue")
	}))
	requireNoError(t, err)
	requireListOrder(t, label, first)

	annotation, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.AnnotationEquals("tier", "backend")
	}))
	requireNoError(t, err)
	requireListOrder(t, annotation, second)
}

func TestQueryReturnsDetachedItems(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 1, listItem(key, 1, "stored"))

	result, err := cache.Query(objectquery.Predicate{})
	requireNoError(t, err)
	mutateState(&result.Items[0].State, "mutated")

	again, err := cache.Query(objectquery.Predicate{})
	requireNoError(t, err)
	if desired := desiredString(t, again.Items[0].State); desired != "stored" {
		t.Fatalf("desired = %q; want stored", desired)
	}
}

func TestQueryIgnoresHistoryAndDeletedTombstones(t *testing.T) {
	key := testKey("system", 1)
	cache := readyHistoryCache(t, 2, 1, listItem(key, 1, "match"))

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 1, "match"), 2),
	))

	result, err := cache.Query(builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "match")
	}))
	requireNoError(t, err)
	if result.Revision != 2 || len(result.Items) != 0 {
		t.Fatalf("Query() = %#v; want no latest matches at revision 2", result)
	}

	historical, err := cache.GetAt(key, 1)
	requireNoError(t, err)
	if !historical.Found {
		t.Fatalf("GetAt() Found = false; want history retained")
	}
}

func TestQueryReflectsLatestCreateUpdateDelete(t *testing.T) {
	key := testKey("system", 1)
	cache := readyCache(t, 0)
	predicate := builtPredicate(t, func() (objectquery.Query, error) {
		return objectquery.LabelEquals("env", "active")
	})

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustCreatedChange(key, testState(key, 1, "inactive")),
	))
	result, err := cache.Query(predicate)
	requireNoError(t, err)
	if len(result.Items) != 0 {
		t.Fatalf("Query() items = %d; want 0 before membership", len(result.Items))
	}

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustUpdatedChange(key, testState(key, 1, "inactive"), testState(key, 2, "active")),
	))
	result, err = cache.Query(predicate)
	requireNoError(t, err)
	requireListOrder(t, result, key)

	requireNoError(t, cache.ApplyChange(
		context.Background(),
		objectstore.MustDeletedChange(key, testState(key, 2, "active"), 3),
	))
	result, err = cache.Query(predicate)
	requireNoError(t, err)
	if result.Revision != 3 || len(result.Items) != 0 {
		t.Fatalf("Query() = %#v; want no latest matches after delete", result)
	}
}

func TestQueryAfterReplaceReflectsReplacementStateAndRevision(t *testing.T) {
	oldKey := testKey("system", 1)
	newKey := testKey("system", 2)
	cache := readyCache(t, 5, listItem(oldKey, 1, "old"))

	requireNoError(t, cache.Replace(
		context.Background(),
		collectionRead(t, testCollection(), 8, listItem(newKey, 2, "new")),
	))

	result, err := cache.Query(objectquery.Predicate{})
	requireNoError(t, err)
	if result.Revision != 8 {
		t.Fatalf("revision = %s; want 8", result.Revision)
	}
	requireListOrder(t, result, newKey)
}

func builtPredicate(t testing.TB, build func() (objectquery.Query, error)) objectquery.Predicate {
	t.Helper()

	query, err := build()
	requireNoError(t, err)

	return mustPredicate(t, query)
}

func mustPredicate(t testing.TB, query objectquery.Query) objectquery.Predicate {
	t.Helper()

	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	return predicate
}

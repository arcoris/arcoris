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
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCacheListNilReceiver(t *testing.T) {
	var cache *Cache

	got, err := cache.List(objectquery.Query{})
	requireNoError(t, err)

	if got.Items != nil {
		t.Fatalf("Items = %#v; want nil", got.Items)
	}
	if !got.Revision.IsZero() {
		t.Fatalf("Revision = %v; want zero", got.Revision)
	}
}

func TestCacheListUsesQueryIndexesAndPredicate(t *testing.T) {
	source := testItems()
	cache := mustCache(t, testListResult(33, source...))
	tests := []struct {
		name  string
		query objectquery.Query
	}{
		{
			name: "identity namespace",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
			},
		},
		{
			name: "identity name",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Name: mustNameEquals(t, "worker-3"),
				},
			},
		},
		{
			name: "label positive negative",
			query: objectquery.Query{
				Labels: mustLabelSelector(
					t,
					mustLabelExists(t, "env"),
					mustLabelNotEquals(t, "tier", "frontend"),
				),
			},
		},
		{
			name: "annotation in",
			query: objectquery.Query{
				Annotations: mustAnnotationSelector(t, mustAnnotationIn(t, "team", "core", "tools")),
			},
		},
		{
			name: "combined",
			query: objectquery.Query{
				Identity: objectquery.IdentitySelector{
					Namespace: mustNamespaceEquals(t, "system"),
				},
				Labels:      mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
				Annotations: mustAnnotationSelector(t, mustAnnotationNotIn(t, "zone", "west")),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertCacheListEquivalent(t, cache, tt.query)
		})
	}
}

func TestCacheListPreservesOrderAndRevision(t *testing.T) {
	source := []objectstore.ListItem{
		testItem("system", "worker-3", 3, labelsMap("env", "prod"), nil),
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 2, labelsMap("env", "qa"), nil),
		testItem("system", "worker-4", 4, labelsMap("env", "prod"), nil),
	}
	cache := mustCache(t, testListResult(34, source...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	got, err := cache.List(query)
	requireNoError(t, err)

	if got.Revision != 34 {
		t.Fatalf("Revision = %v; want 34", got.Revision)
	}
	requireItemOrder(
		t,
		got.Items,
		itemRef{"system", "worker-3", 3},
		itemRef{"system", "worker-1", 1},
		itemRef{"system", "worker-4", 4},
	)
}

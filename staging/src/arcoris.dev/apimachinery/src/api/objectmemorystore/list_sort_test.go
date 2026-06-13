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

package objectmemorystore

import (
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestSortListItemsOrdersByFullStorageKey(t *testing.T) {
	want := []objectstore.Key{
		sortKey("apps.arcoris.dev", "v1", "workers", "alpha", "one"),
		sortKey("control.arcoris.dev", "v1", "jobs", "alpha", "one"),
		sortKey("control.arcoris.dev", "v1", "workers", "alpha", "one"),
		sortKey("control.arcoris.dev", "v1", "workers", "alpha", "two"),
		sortKey("control.arcoris.dev", "v1", "workers", "beta", "one"),
		sortKey("control.arcoris.dev", "v2", "workers", "alpha", "one"),
	}
	items := []objectstore.ListItem{
		{Key: want[5]},
		{Key: want[2]},
		{Key: want[4]},
		{Key: want[0]},
		{Key: want[3]},
		{Key: want[1]},
	}

	sortListItems(items)

	for i, item := range items {
		if !item.Key.Equal(want[i]) {
			t.Fatalf("item[%d] = %s; want %s", i, item.Key, want[i])
		}
	}
}

func TestCompareListKeysReportsEqualKeys(t *testing.T) {
	key := sortKey("control.arcoris.dev", "v1", "workers", "system", "main")

	if got := compareListKeys(key, key); got != 0 {
		t.Fatalf("compareListKeys(equal) = %d; want 0", got)
	}
}

// sortKey constructs a validated storage key for list ordering tests.
func sortKey(group, version, resourceName, namespace, name string) objectstore.Key {
	return objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    apiidentity.Group(group),
			Version:  apiidentity.Version(version),
			Resource: apiidentity.Resource(resourceName),
		},
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace(namespace),
			Name:      metaidentity.Name(name),
		},
	)
}

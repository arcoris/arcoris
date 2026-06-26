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

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestCollectionCloneDetachesOrderMapAndStates(t *testing.T) {
	key := testKey("system", 1)
	col, err := buildCollection(objectstore.ListResult{
		Items:    []objectstore.ListItem{listItem(key, 1, "cached")},
		Revision: 1,
	})
	requireNoError(t, err)

	clone := col.clone()
	clone.order[0] = testKey("system", 2)
	item := clone.items[key]
	mutateState(&item.State, "mutated")
	clone.items[key] = item

	got := col.listResult()
	requireListOrder(t, got, key)
	if desired := desiredString(t, got.Items[0].State); desired != "cached" {
		t.Fatalf("clone mutation changed original desired to %q", desired)
	}
}

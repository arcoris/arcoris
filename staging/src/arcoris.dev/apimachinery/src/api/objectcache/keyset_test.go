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

import "testing"

func TestKeySetOperations(t *testing.T) {
	items := testItems()
	left := newKeySet(items[0].Key, items[1].Key)
	right := newKeySet(items[1].Key, items[2].Key)

	union := unionKeySets(left, right)
	for _, item := range items[:3] {
		if !union.has(item.Key) {
			t.Fatalf("union missing %s", item.Key)
		}
	}

	intersection := intersectKeySets(left, right)
	if !intersection.has(items[1].Key) {
		t.Fatalf("intersection missing %s", items[1].Key)
	}
	if intersection.has(items[0].Key) || intersection.has(items[2].Key) {
		t.Fatalf("intersection = %#v; want only worker-2", intersection)
	}
}

func TestKeySetCloneDetached(t *testing.T) {
	item := testItems()[0]
	original := newKeySet(item.Key)
	cloned := original.clone()

	cloned.remove(item.Key)

	if !original.has(item.Key) {
		t.Fatal("mutating clone removed key from original")
	}
}

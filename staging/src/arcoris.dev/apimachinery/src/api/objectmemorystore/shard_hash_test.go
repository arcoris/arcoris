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

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestHashKeyIsDeterministic(t *testing.T) {
	key := testKey(1)

	if first, second := hashKey(key), hashKey(key); first != second {
		t.Fatalf("hashKey changed for same key: %d != %d", first, second)
	}
}

func TestHashKeyChangesWithKeyComponents(t *testing.T) {
	base := testKey(1)
	changedResource := base
	changedResource.Resource.Resource = "controllers"
	changedNamespace := base
	changedNamespace.Object.Namespace = "other"
	changedName := testKey(2)

	tests := []struct {
		name string
		key  objectKey
	}{
		{name: "resource", key: objectKey{base: base, changed: changedResource}},
		{name: "namespace", key: objectKey{base: base, changed: changedNamespace}},
		{name: "name", key: objectKey{base: base, changed: changedName}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if hashKey(tt.key.base) == hashKey(tt.key.changed) {
				t.Fatalf("hashKey did not change when %s changed", tt.name)
			}
		})
	}
}

func TestHashKeyAllocatesNoMemory(t *testing.T) {
	key := testKey(1)

	allocs := testing.AllocsPerRun(1000, func() {
		_ = hashKey(key)
	})

	if allocs != 0 {
		t.Fatalf("allocs = %v; want 0", allocs)
	}
}

func TestHashKeyIsIndependentOfMapIterationOrder(t *testing.T) {
	keys := map[objectstore.Key]uint64{
		testKey(1): hashKey(testKey(1)),
		testKey(2): hashKey(testKey(2)),
		testKey(3): hashKey(testKey(3)),
	}

	for range 10 {
		for key, want := range keys {
			if got := hashKey(key); got != want {
				t.Fatalf("hashKey(%s) = %d; want %d", key, got, want)
			}
		}
	}
}

type objectKey struct {
	base    objectstore.Key
	changed objectstore.Key
}

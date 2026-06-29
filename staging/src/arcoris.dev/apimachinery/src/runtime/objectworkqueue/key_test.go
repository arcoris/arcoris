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

package objectworkqueue

import "testing"

func TestObjectStoreKeyCanBeUsedAsMapIdentity(t *testing.T) {
	key := testKey(1)
	seen := map[keyID]struct{}{
		key: {},
	}

	if _, ok := seen[key]; !ok {
		t.Fatalf("key lookup failed")
	}
}

func TestKeyForItemReturnsObjectStoreKey(t *testing.T) {
	item := testItem(1)

	if !keyForItem(item).Equal(item.Key) {
		t.Fatalf("keyForItem() = %s; want %s", keyForItem(item), item.Key)
	}
}

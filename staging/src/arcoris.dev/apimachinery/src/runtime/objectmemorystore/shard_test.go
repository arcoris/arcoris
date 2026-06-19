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

import "testing"

func TestShardGetOrCreateReturnsStableSlot(t *testing.T) {
	var shard shard
	shard.init()
	key := testKey(1)

	first := shard.getOrCreate(key)
	second := shard.getOrCreate(key)

	if first == nil {
		t.Fatalf("getOrCreate returned nil")
	}
	if first != second {
		t.Fatalf("getOrCreate returned different slots for same key")
	}
	if got := shard.get(key); got != first {
		t.Fatalf("get returned %p; want %p", got, first)
	}
}

func TestStoreShardForIsDeterministic(t *testing.T) {
	store := testStore(t, WithShardCount(4))
	key := testKey(1)

	if store.shardFor(key) != store.shardFor(key) {
		t.Fatalf("shardFor changed for same key")
	}
}

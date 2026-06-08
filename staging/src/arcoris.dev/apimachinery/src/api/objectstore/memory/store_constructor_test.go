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

package memory

import "testing"

func TestNewBuildsInitializedStore(t *testing.T) {
	store, err := New()
	requireNoError(t, err)

	if store == nil {
		t.Fatalf("New() returned nil store")
	}
	if len(store.shards) != int(defaultShardCount) {
		t.Fatalf("shards = %d; want %d", len(store.shards), defaultShardCount)
	}
	if store.mask != uint64(defaultShardCount-1) {
		t.Fatalf("mask = %d; want %d", store.mask, defaultShardCount-1)
	}
	for i := range store.shards {
		if store.shards[i].slots == nil {
			t.Fatalf("shard %d was not initialized", i)
		}
	}
}

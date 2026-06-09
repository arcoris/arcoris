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

// New constructs an empty in-memory object store.
func New(opts ...Option) (*Store, error) {
	cfg, err := applyOptions(opts)
	if err != nil {
		return nil, err
	}

	store := &Store{
		shards: make([]shard, cfg.shardCount),
		mask:   uint64(cfg.shardCount - 1),
	}
	for i := range store.shards {
		store.shards[i].init()
	}

	return store, nil
}

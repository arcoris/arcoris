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

// Option configures a new memory object store.
type Option func(*config)

// WithShardCount sets the number of fixed key-index shards.
//
// The value must be a non-zero power of two. The power-of-two restriction keeps
// shard selection cheap and deterministic without modulo bias details leaking
// into the rest of the implementation.
func WithShardCount(n uint) Option {
	return func(cfg *config) {
		cfg.shardCount = n
	}
}

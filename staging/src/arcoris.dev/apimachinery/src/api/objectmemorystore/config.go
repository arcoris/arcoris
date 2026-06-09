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

// defaultShardCount is the fixed key-index width used when callers do not
// configure a store.
const defaultShardCount = uint(64)

// config stores normalized construction parameters.
type config struct {
	// shardCount is fixed after New returns and never changes at runtime.
	shardCount uint
}

// defaultConfig returns the stable package defaults.
func defaultConfig() config {
	return config{shardCount: defaultShardCount}
}

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

// Snapshot returns an immutable detached view of the current cache state.
func (c *Cache) Snapshot() Snapshot {
	if c == nil {
		return Snapshot{}
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return Snapshot{col: c.col.clone()}
}

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

import "arcoris.dev/apimachinery/api/objectstore"

// Revision returns the current collection boundary revision and readiness.
func (c *Cache) Revision() (objectstore.Revision, bool) {
	if c == nil {
		return 0, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.ready {
		return 0, false
	}

	return c.col.revision, true
}

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

package typecatalog

import "arcoris.dev/apimachinery/api/types"

// Definitions returns registered definitions in stable registration order.
//
// Definition has private descriptor payload and returns detached Descriptor
// values from Descriptor(), so the returned values do not expose mutable catalog
// internals.
func (c *Catalog) Definitions() []types.Definition {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]types.Definition, 0, len(c.order))
	for _, name := range c.order {
		out = append(out, c.defs[name])
	}
	return out
}

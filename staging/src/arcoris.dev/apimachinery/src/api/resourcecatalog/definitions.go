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

package resourcecatalog

import (
	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

// Definitions returns registered definitions in stable registration order.
//
// The returned slice is detached. Definition keeps its version slice private
// and returns detached slices from Versions(), so returning descriptor values
// does not expose mutable catalog storage.
func (c *Catalog) Definitions() []resource.Definition {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make([]resource.Definition, 0, len(c.order))
	for _, gr := range c.order {
		out = append(out, c.defsByResource[gr])
	}
	return out
}

// Resources returns registered GroupResource keys in stable registration order.
//
// The returned slice is detached. Mutating it does not affect the catalog.
func (c *Catalog) Resources() []identity.GroupResource {
	if c == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return append([]identity.GroupResource(nil), c.order...)
}

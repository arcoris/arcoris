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

// cloneLocked returns a detached catalog copy while c.mu is already held.
//
// RegisterMany uses the clone as a candidate resolver before mutating the
// receiver. That gives atomic registration its key invariant: validation can
// see existing definitions and same-batch definitions, but failed validation
// cannot leak partial state into the owner catalog.
func (c *Catalog) cloneLocked() *Catalog {
	out := &Catalog{
		defs:  make(map[types.TypeName]types.TypeDefinition, len(c.defs)),
		order: append([]types.TypeName(nil), c.order...),
	}
	for name, def := range c.defs {
		out.defs[name] = def
	}
	return out
}

// storeLocked stores def while the caller holds the catalog lock.
//
// The helper owns lazy map initialization and registration-order maintenance so
// write paths cannot accidentally update one without the other.
func (c *Catalog) storeLocked(def types.TypeDefinition) {
	if c.defs == nil {
		c.defs = make(map[types.TypeName]types.TypeDefinition)
	}
	c.defs[def.Name()] = def
	c.order = append(c.order, def.Name())
}

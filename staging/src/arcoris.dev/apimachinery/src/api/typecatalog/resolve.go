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

// Resolve returns the definition registered for name.
//
// A nil catalog behaves like an empty catalog. Definition keeps its
// descriptor payload private, and Definition.Descriptor returns a detached Descriptor,
// so returning the value is safe for callers outside package types.
//
// Typical resolver use:
//
//	desc := types.Ref("meta.arcoris.dev.Name").
//		Descriptor()
//
//	err := types.ValidateResolved(
//		desc,
//		&catalog,
//	)
func (c *Catalog) Resolve(name types.TypeName) (types.Definition, bool) {
	if c == nil {
		return types.Definition{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	def, ok := c.defs[name]
	if !ok {
		return types.Definition{}, false
	}
	return def, true
}

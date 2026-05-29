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

import (
	"fmt"
	"sync"

	"arcoris.dev/apimachinery/api/types"
)

// Catalog owns named structural type definitions.
//
// Catalog is zero-value usable, concurrency-safe, and explicitly owner-created.
// It deliberately has no package-level singleton or init-time registration
// path. That keeps descriptor composition separate from runtime schemes,
// resource registries, codec registries, and process-wide extension state.
type Catalog struct {
	// mu protects defs and order.
	//
	// Catalog supports concurrent readers and serialized writers. Registration
	// is intentionally serialized because batch validation must observe a stable
	// candidate catalog.
	mu sync.RWMutex

	// defs stores definitions by validated name.
	//
	// The map is initialized lazily so the zero value remains usable.
	defs map[types.TypeName]types.TypeDefinition

	// order preserves stable registration order for enumeration.
	//
	// Catalog enumeration follows owner registration order rather than Go map
	// iteration order.
	order []types.TypeName
}

// Register validates and stores one definition.
//
// Register is a convenience wrapper around RegisterMany. It has the same
// validation, duplicate-rejection, and atomicity behavior.
//
// Typical registration:
//
//	var catalog Catalog
//
//	err := catalog.Register(types.Define(
//		"arcoris.meta.Name",
//		types.String().
//			MinLen(1).
//			MaxLen(253),
//	))
func (c *Catalog) Register(def types.TypeDefinition) error {
	return c.RegisterMany(def)
}

// RegisterMany validates and stores definitions atomically.
//
// If any definition is invalid, duplicated within the batch, or conflicts with
// an existing catalog name, the receiver remains unchanged. Validation runs
// against a cloned candidate catalog containing existing definitions plus the
// entire incoming batch, so definitions may reference siblings registered in
// the same call.
//
// Typical batch registration:
//
//	err := catalog.RegisterMany(
//		types.Define(
//			"arcoris.meta.Name",
//			types.String().
//				MinLen(1),
//		),
//		types.Define(
//			"arcoris.meta.NameList",
//			types.ListOf(
//				types.Ref("arcoris.meta.Name"),
//			),
//		),
//	)
func (c *Catalog) RegisterMany(defs ...types.TypeDefinition) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	seen := make(map[types.TypeName]struct{}, len(defs))
	for i, def := range defs {
		name := def.Name()
		if !name.IsValid() {
			return catalogError(fmt.Sprintf("definitions[%d].name", i), types.ErrInvalidTypeReference)
		}
		// Duplicate names are catalog ownership conflicts. They are not TypeRef
		// syntax or resolution failures, so they keep catalog-specific error
		// identity instead of joining types.ErrInvalidTypeReference.
		if _, ok := seen[name]; ok {
			return catalogError(fmt.Sprintf("definitions[%d].name", i), ErrDuplicateDefinition)
		}
		seen[name] = struct{}{}
		if _, ok := c.defs[name]; ok {
			return catalogError(fmt.Sprintf("definitions[%s]", name), ErrDefinitionExists)
		}
	}

	candidate := c.cloneLocked()
	for _, def := range defs {
		candidate.storeLocked(def)
	}
	for i, def := range defs {
		if err := types.ValidateDefinition(def, candidate); err != nil {
			return catalogError(fmt.Sprintf("definitions[%d]", i), err)
		}
	}

	if c.defs == nil {
		c.defs = make(map[types.TypeName]types.TypeDefinition, len(defs))
	}
	for _, def := range defs {
		c.storeLocked(def)
	}
	return nil
}

// ResolveType returns the definition registered for name.
//
// A nil catalog behaves like an empty catalog. TypeDefinition keeps its
// descriptor payload private, and TypeDefinition.Type returns a detached Type,
// so returning the value is safe for callers outside package types.
//
// Typical resolver use:
//
//	tp := types.Ref("arcoris.meta.Name").
//		Type()
//
//	err := types.ValidateType(
//		tp,
//		&catalog,
//	)
func (c *Catalog) ResolveType(name types.TypeName) (types.TypeDefinition, bool) {
	if c == nil {
		return types.TypeDefinition{}, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	def, ok := c.defs[name]
	if !ok {
		return types.TypeDefinition{}, false
	}
	return def, true
}

// Names returns registered names in stable registration order.
//
// The returned slice is detached. Mutating it does not affect the catalog.
func (c *Catalog) Names() []types.TypeName {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]types.TypeName(nil), c.order...)
}

// Definitions returns registered definitions in stable registration order.
//
// TypeDefinition has private descriptor payload and returns detached Type
// values from Type(), so the returned values do not expose mutable catalog
// internals.
func (c *Catalog) Definitions() []types.TypeDefinition {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]types.TypeDefinition, 0, len(c.order))
	for _, name := range c.order {
		out = append(out, c.defs[name])
	}
	return out
}

// cloneLocked returns a detached catalog copy while c.mu is already held.
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
func (c *Catalog) storeLocked(def types.TypeDefinition) {
	if c.defs == nil {
		c.defs = make(map[types.TypeName]types.TypeDefinition)
	}
	c.defs[def.Name()] = def
	c.order = append(c.order, def.Name())
}

// catalogError adds catalog location context while preserving error identity.
func catalogError(path string, err error) error {
	return fmt.Errorf("typecatalog: %s: %w", path, err)
}

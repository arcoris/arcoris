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

	"arcoris.dev/apimachinery/api/types"
)

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
//			MinBytes(1).
//			MaxRunes(253),
//	))
func (c *Catalog) Register(def types.Definition) error {
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
//				MinBytes(1),
//		),
//		types.Define(
//			"arcoris.meta.NameList",
//			types.ListOf(
//				types.Ref("arcoris.meta.Name"),
//			),
//		),
//	)
func (c *Catalog) RegisterMany(defs ...types.Definition) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.checkIncomingDefinitionsLocked(defs); err != nil {
		return err
	}

	candidate := c.cloneLocked()
	for _, def := range defs {
		candidate.storeLocked(def)
	}
	for i, def := range defs {
		if err := types.ValidateDefinitionResolved(def, candidate); err != nil {
			return catalogError(fmt.Sprintf("definitions[%d]", i), err)
		}
	}

	for _, def := range defs {
		c.storeLocked(def)
	}
	return nil
}

// checkIncomingDefinitionsLocked checks batch names before mutating state.
//
// Duplicate names are catalog ownership conflicts. They are not DescriptorRef syntax
// or resolution failures, so they keep catalog-specific error identity instead
// of joining types.ErrInvalidDescriptorReference.
func (c *Catalog) checkIncomingDefinitionsLocked(defs []types.Definition) error {
	seen := make(map[types.TypeName]struct{}, len(defs))
	for i, def := range defs {
		name := def.Name()
		if !name.IsValid() {
			return catalogError(
				fmt.Sprintf("definitions[%d].name", i),
				types.ErrInvalidDescriptorReference,
			)
		}

		if _, ok := seen[name]; ok {
			return catalogError(
				fmt.Sprintf("definitions[%d].name", i),
				ErrDuplicateDefinition,
			)
		}
		seen[name] = struct{}{}

		if _, ok := c.defs[name]; ok {
			return catalogError(
				fmt.Sprintf("definitions[%s]", name),
				ErrDefinitionExists,
			)
		}
	}
	return nil
}

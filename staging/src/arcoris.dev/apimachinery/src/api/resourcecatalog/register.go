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
	"fmt"

	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
)

// New creates an owner-created resource catalog.
//
// The resolver is used only during registration to validate resource surface
// DescriptorRef roots. It is not mutable after construction because changing it would
// weaken the catalog's registration-time invariants.
func New(resolver types.Resolver) *Catalog {
	return &Catalog{resolver: resolver}
}

// Register validates and stores one resource definition.
//
// Register is a convenience wrapper around RegisterMany. It has the same
// validation, duplicate-rejection, conflict-rejection, and atomicity behavior.
func (c *Catalog) Register(def resource.Definition) error {
	return c.RegisterMany(def)
}

// RegisterMany validates and stores definitions atomically.
//
// If any definition is invalid, duplicated within the batch, or conflicts with
// existing catalog state, the receiver remains unchanged. Validation runs
// against a cloned candidate catalog containing existing definitions plus the
// entire incoming batch.
func (c *Catalog) RegisterMany(defs ...resource.Definition) error {
	if c == nil {
		return nilCatalogError()
	}

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
		if err := validateCatalogDefinition(def, candidate.resolver); err != nil {
			return nestedCatalogError(
				definitionPath(i),
				ErrorReasonInvalidDefinition,
				fmt.Sprintf("resource definition %d is invalid", i),
				err,
			)
		}
	}

	for _, def := range defs {
		c.storeLocked(def)
	}
	return nil
}

// validateCatalogDefinition applies the catalog's registration-time validation
// policy to one resource definition.
//
// A catalog with a resolver can prove DescriptorRef roots through resolved
// validation. A catalog without a resolver can still accept direct-object
// definitions through local validation, but rejects root refs because they
// cannot be proven object-like at registration time.
func validateCatalogDefinition(def resource.Definition, resolver types.Resolver) error {
	if resolver != nil {
		return def.ValidateResolved(resolver)
	}

	if err := def.ValidateLocal(); err != nil {
		return err
	}

	for _, version := range def.Versions() {
		if version.Desired().Code() == types.DescriptorRef {
			return fmt.Errorf(
				"definition.versions[%s].desired requires a catalog resolver: %w",
				version.Version(),
				resource.ErrInvalidVersion,
			)
		}
		if observed, ok := version.Observed(); ok && observed.Code() == types.DescriptorRef {
			return fmt.Errorf(
				"definition.versions[%s].observed requires a catalog resolver: %w",
				version.Version(),
				resource.ErrInvalidVersion,
			)
		}
	}

	return nil
}

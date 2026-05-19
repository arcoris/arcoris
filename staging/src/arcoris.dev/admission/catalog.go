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

package admission

// nilCatalogPanic is the stable panic string for nil *Catalog method receivers.
//
// A nil Catalog is a wiring/configuration bug. Missing catalog entries are
// represented by lookup methods returning false.
const nilCatalogPanic = "admission.Catalog: nil catalog"

// Catalog aggregates owner-created admission metadata registries.
//
// Catalog is a convenience boundary for reason, kind, and component catalog
// access. It is not global, not a runtime instance registry, does not store live
// Admitter values, and does not execute admission chains. Each method delegates
// to the underlying concurrency-safe registries.
type Catalog struct {
	// reasons stores stable reason descriptors for documentation, config
	// validation, and higher-level catalog checks.
	reasons *ReasonRegistry

	// kinds stores stable component kind descriptors.
	kinds *KindRegistry

	// components stores stable component descriptors validated against kinds.
	components *ComponentRegistry
}

// NewCatalog creates an aggregate catalog from owner-provided registries.
//
// All registries are required. Passing nil would create ambiguous partial
// catalog behavior, so construction rejects nil registries explicitly.
func NewCatalog(
	reasons *ReasonRegistry,
	kinds *KindRegistry,
	components *ComponentRegistry,
) (*Catalog, error) {
	if reasons == nil {
		return nil, ErrNilReasonRegistry
	}
	if kinds == nil {
		return nil, ErrNilKindRegistry
	}
	if components == nil {
		return nil, ErrNilComponentRegistry
	}

	return &Catalog{
		reasons:    reasons,
		kinds:      kinds,
		components: components,
	}, nil
}

// Reason returns the descriptor registered for reason.
func (c *Catalog) Reason(reason Reason) (ReasonDescriptor, bool) {
	c.requireNonNil()
	return c.reasons.Lookup(reason)
}

// Kind returns the descriptor registered for kind.
func (c *Catalog) Kind(kind ComponentKind) (ComponentKindDescriptor, bool) {
	c.requireNonNil()
	return c.kinds.Lookup(kind)
}

// Component returns the descriptor registered for id.
func (c *Catalog) Component(id ComponentID) (ComponentDescriptor, bool) {
	c.requireNonNil()
	return c.components.Lookup(id)
}

// Reasons returns sorted, copy-safe reason descriptors.
func (c *Catalog) Reasons() []ReasonDescriptor {
	c.requireNonNil()
	return c.reasons.List()
}

// Kinds returns sorted, copy-safe component kind descriptors.
func (c *Catalog) Kinds() []ComponentKindDescriptor {
	c.requireNonNil()
	return c.kinds.List()
}

// Components returns sorted, copy-safe component descriptors.
func (c *Catalog) Components() []ComponentDescriptor {
	c.requireNonNil()
	return c.components.List()
}

// RegisterReason adds descriptor to the catalog's ReasonRegistry.
func (c *Catalog) RegisterReason(descriptor ReasonDescriptor) error {
	c.requireNonNil()
	return c.reasons.Register(descriptor)
}

// RegisterKind adds descriptor to the catalog's KindRegistry.
func (c *Catalog) RegisterKind(descriptor ComponentKindDescriptor) error {
	c.requireNonNil()
	return c.kinds.Register(descriptor)
}

// RegisterComponent adds descriptor to the catalog's ComponentRegistry.
//
// Component registration is still validated through ComponentRegistry.Register.
// In the normal catalog shape the component registry is constructed with the
// same KindRegistry passed to NewCatalog, so custom kinds registered through the
// catalog are available to later component registrations.
func (c *Catalog) RegisterComponent(descriptor ComponentDescriptor) error {
	c.requireNonNil()
	return c.components.Register(descriptor)
}

// LenReasons reports the number of registered reasons.
func (c *Catalog) LenReasons() int {
	c.requireNonNil()
	return c.reasons.Len()
}

// LenKinds reports the number of registered component kinds.
func (c *Catalog) LenKinds() int {
	c.requireNonNil()
	return c.kinds.Len()
}

// LenComponents reports the number of registered components.
func (c *Catalog) LenComponents() int {
	c.requireNonNil()
	return c.components.Len()
}

// requireNonNil panics with the stable nil receiver message used by all Catalog
// methods.
func (c *Catalog) requireNonNil() {
	if c == nil {
		panic(nilCatalogPanic)
	}
}

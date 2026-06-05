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

package admissioncatalog

import "arcoris.dev/admission"

// Reason returns the descriptor declared for reason.
//
// Invalid reasons and missing valid reasons both return false. Lookup is not a
// validation report.
func (c *Catalog) Reason(reason admission.Reason) (ReasonDescriptor, bool) {
	c.requireNonNil()
	if !reason.IsValid() {
		return ReasonDescriptor{}, false
	}
	return c.reasons.get(reason)
}

// HasReason reports whether reason has been declared in the catalog.
func (c *Catalog) HasReason(reason admission.Reason) bool {
	_, ok := c.Reason(reason)
	return ok
}

// Reasons returns declared reason descriptors sorted by reason string.
//
// The returned slice is a fresh copy.
func (c *Catalog) Reasons() []ReasonDescriptor {
	c.requireNonNil()
	return c.reasons.list()
}

// LenReasons reports the number of declared reasons.
func (c *Catalog) LenReasons() int {
	c.requireNonNil()
	return c.reasons.len()
}

// Kind returns the descriptor declared for kind.
//
// Invalid kinds and missing valid kinds both return false. Lookup is not a
// validation report.
func (c *Catalog) Kind(kind admission.ComponentKind) (ComponentKindDescriptor, bool) {
	c.requireNonNil()
	if !kind.IsValid() {
		return ComponentKindDescriptor{}, false
	}
	return c.kinds.get(kind)
}

// HasKind reports whether kind has been declared in the catalog.
func (c *Catalog) HasKind(kind admission.ComponentKind) bool {
	_, ok := c.Kind(kind)
	return ok
}

// Kinds returns declared component kind descriptors sorted by kind string.
//
// The returned slice is a fresh copy.
func (c *Catalog) Kinds() []ComponentKindDescriptor {
	c.requireNonNil()
	return c.kinds.list()
}

// LenKinds reports the number of declared component kinds.
func (c *Catalog) LenKinds() int {
	c.requireNonNil()
	return c.kinds.len()
}

// Component returns the descriptor declared for id.
//
// Invalid IDs and missing valid IDs both return false. Lookup is not a
// validation report.
func (c *Catalog) Component(id admission.ComponentID) (ComponentDescriptor, bool) {
	c.requireNonNil()
	if !id.IsValid() {
		return ComponentDescriptor{}, false
	}
	return c.components.get(id)
}

// HasComponent reports whether id has been declared in the catalog.
func (c *Catalog) HasComponent(id admission.ComponentID) bool {
	_, ok := c.Component(id)
	return ok
}

// Components returns declared component descriptors sorted by component ID.
//
// The returned slice is a fresh copy.
func (c *Catalog) Components() []ComponentDescriptor {
	c.requireNonNil()
	return c.components.list()
}

// LenComponents reports the number of declared components.
func (c *Catalog) LenComponents() int {
	c.requireNonNil()
	return c.components.len()
}

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

// Reason returns the descriptor registered for reason.
func (c *Catalog) Reason(reason admission.Reason) (ReasonDescriptor, bool) {
	c.requireNonNil()
	return c.reasons.Lookup(reason)
}

// Kind returns the descriptor registered for kind.
func (c *Catalog) Kind(kind admission.ComponentKind) (ComponentKindDescriptor, bool) {
	c.requireNonNil()
	return c.kinds.Lookup(kind)
}

// Component returns the descriptor registered for id.
func (c *Catalog) Component(id admission.ComponentID) (ComponentDescriptor, bool) {
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

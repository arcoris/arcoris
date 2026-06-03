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

// RegisterReason adds descriptor to the catalog's ReasonRegistry.
func (c *Catalog) RegisterReason(descriptor ReasonDescriptor) error {
	c.requireNonNil()
	return c.reasons.Register(descriptor)
}

// RegisterKind adds descriptor to the kind catalog used by later component
// registration.
//
// NewCatalog guarantees that the ComponentRegistry delegates validation to this
// same KindRegistry reference, so a newly registered kind can be used by a
// later RegisterComponent call.
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

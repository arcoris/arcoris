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

// ensureStorageLocked initializes catalog maps while the caller holds c.mu.
func (c *Catalog) ensureStorageLocked() {
	if c.defsByResource == nil {
		c.defsByResource = make(map[identity.GroupResource]resource.Definition)
	}
	if c.resourceByKind == nil {
		c.resourceByKind = make(map[identity.GroupKind]identity.GroupResource)
	}
	if c.versionByResource == nil {
		c.versionByResource = make(map[identity.GroupVersionResource]versionRef)
	}
	if c.versionByKind == nil {
		c.versionByKind = make(map[identity.GroupVersionKind]versionRef)
	}
}

// cloneLocked returns a detached candidate catalog while c.mu is already held.
//
// RegisterMany validates against the clone before mutating the receiver. That
// preserves atomicity while allowing validation to see existing definitions and
// same-batch definitions.
func (c *Catalog) cloneLocked() *Catalog {
	out := &Catalog{
		resolver:          c.resolver,
		defsByResource:    make(map[identity.GroupResource]resource.Definition, len(c.defsByResource)),
		order:             append([]identity.GroupResource(nil), c.order...),
		resourceByKind:    make(map[identity.GroupKind]identity.GroupResource, len(c.resourceByKind)),
		versionByResource: make(map[identity.GroupVersionResource]versionRef, len(c.versionByResource)),
		versionByKind:     make(map[identity.GroupVersionKind]versionRef, len(c.versionByKind)),
	}

	for key, def := range c.defsByResource {
		out.defsByResource[key] = def
	}
	for key, value := range c.resourceByKind {
		out.resourceByKind[key] = value
	}
	for key, value := range c.versionByResource {
		out.versionByResource[key] = value
	}
	for key, value := range c.versionByKind {
		out.versionByKind[key] = value
	}

	return out
}

// storeLocked stores def and all derived indexes.
//
// The caller must already hold the catalog lock and must already have validated
// and conflict-checked the definition. The helper owns map initialization and
// index population only.
func (c *Catalog) storeLocked(def resource.Definition) {
	c.ensureStorageLocked()

	gr := groupResourceOf(def)
	gk := groupKindOf(def)

	c.defsByResource[gr] = def
	c.order = append(c.order, gr)
	c.resourceByKind[gk] = gr

	for _, gvr := range versionResourceKeys(def) {
		c.versionByResource[gvr] = versionRef{
			resource: gr,
			version:  gvr.Version,
		}
	}
	for _, gvk := range versionKindKeys(def) {
		c.versionByKind[gvk] = versionRef{
			resource: gr,
			version:  gvk.Version,
		}
	}
}

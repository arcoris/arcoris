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

// ResolveResource returns the definition registered for gr.
//
// A nil catalog behaves like an empty catalog. Resolve methods use exact map
// lookup and do not validate the input identity.
func (c *Catalog) ResolveResource(gr identity.GroupResource) (resource.Definition, bool) {
	if c == nil {
		return resource.Definition{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	def, ok := c.defsByResource[gr]
	if !ok {
		return resource.Definition{}, false
	}
	return def, true
}

// ResolveKind returns the definition registered for gk.
func (c *Catalog) ResolveKind(gk identity.GroupKind) (resource.Definition, bool) {
	if c == nil {
		return resource.Definition{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	gr, ok := c.resourceByKind[gk]
	if !ok {
		return resource.Definition{}, false
	}

	def, ok := c.defsByResource[gr]
	if !ok {
		return resource.Definition{}, false
	}
	return def, true
}

// ResolveVersionResource returns the definition and exact version registered
// for gvr.
func (c *Catalog) ResolveVersionResource(
	gvr identity.GroupVersionResource,
) (resource.Definition, resource.VersionDefinition, bool) {
	if c == nil {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	ref, ok := c.versionByResource[gvr]
	if !ok {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	return c.resolveVersionLocked(ref)
}

// ResolveVersionKind returns the definition and exact version registered for
// gvk.
func (c *Catalog) ResolveVersionKind(
	gvk identity.GroupVersionKind,
) (resource.Definition, resource.VersionDefinition, bool) {
	if c == nil {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	ref, ok := c.versionByKind[gvk]
	if !ok {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	return c.resolveVersionLocked(ref)
}

// resolveVersionLocked follows a versionRef back to its Definition and exact
// VersionDefinition while the caller holds c.mu for reading.
func (c *Catalog) resolveVersionLocked(
	ref versionRef,
) (resource.Definition, resource.VersionDefinition, bool) {
	def, ok := c.defsByResource[ref.resource]
	if !ok {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	version, ok := def.Version(ref.version)
	if !ok {
		return resource.Definition{}, resource.VersionDefinition{}, false
	}

	return def, version, true
}

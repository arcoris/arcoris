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

package resource

import (
	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/types"
)

// Definition describes one durable API resource family.
//
// A resource family fixes group, kind, resource collection name, scope, and the
// set of versioned API contracts. GVK and GVR values are derived from the family
// identity plus a concrete version.
type Definition struct {
	// group is the API group shared by every version of the resource family.
	group identity.Group

	// kind is the API object kind shared by every version of the resource family.
	kind identity.Kind

	// resource is the collection identity shared by every version of the family.
	resource identity.Resource

	// scope describes future instance addressing without defining metadata.
	scope Scope

	// versions preserves declaration order for version descriptors.
	versions []VersionDefinition
}

// NewDefinition constructs a resource-family definition.
//
// The constructor records the supplied values without validation. Call Validate
// or ValidateDefinition at a trust boundary.
func NewDefinition(
	group identity.Group,
	kind identity.Kind,
	resource identity.Resource,
	scope Scope,
	versions ...VersionDefinition,
) Definition {
	return Definition{
		group:    group,
		kind:     kind,
		resource: resource,
		scope:    scope,
		versions: cloneVersions(versions),
	}
}

// IsZero reports whether d is the empty optional sentinel.
//
// A zero Definition is invalid as a complete resource definition. The method is
// useful for callers that keep optional descriptors in maps or lookup results.
func (d Definition) IsZero() bool {
	return d.group.IsZero() &&
		d.kind.IsZero() &&
		d.resource.IsZero() &&
		d.scope.IsZero() &&
		len(d.versions) == 0
}

// Group returns the API group for the resource family.
func (d Definition) Group() identity.Group { return d.group }

// Kind returns the kind for the resource family.
func (d Definition) Kind() identity.Kind { return d.kind }

// Resource returns the resource collection name for the resource family.
func (d Definition) Resource() identity.Resource { return d.resource }

// Scope returns the future instance-addressing scope for the resource family.
func (d Definition) Scope() Scope { return d.scope }

// Versions returns a detached copy of the version descriptors.
func (d Definition) Versions() []VersionDefinition { return cloneVersions(d.versions) }

// Version returns the descriptor for version.
func (d Definition) Version(version identity.Version) (VersionDefinition, bool) {
	for _, candidate := range d.versions {
		if candidate.version == version {
			return candidate, true
		}
	}
	return VersionDefinition{}, false
}

// GroupKind returns the version-independent kind identity for the family.
func (d Definition) GroupKind() identity.GroupKind {
	return identity.GroupKind{Group: d.group, Kind: d.kind}
}

// GroupResource returns the version-independent resource identity for the family.
func (d Definition) GroupResource() identity.GroupResource {
	return identity.GroupResource{Group: d.group, Resource: d.resource}
}

// GroupVersionKind returns the concrete GVK for version when it is defined.
func (d Definition) GroupVersionKind(version identity.Version) (identity.GroupVersionKind, bool) {
	if _, ok := d.Version(version); !ok {
		return identity.GroupVersionKind{}, false
	}
	return identity.GroupVersionKind{Group: d.group, Version: version, Kind: d.kind}, true
}

// GroupVersionResource returns the concrete GVR for version when it is defined.
func (d Definition) GroupVersionResource(version identity.Version) (identity.GroupVersionResource, bool) {
	if _, ok := d.Version(version); !ok {
		return identity.GroupVersionResource{}, false
	}
	return identity.GroupVersionResource{Group: d.group, Version: version, Resource: d.resource}, true
}

// Validate checks the structural integrity of d.
func (d Definition) Validate(resolver types.Resolver) error {
	return ValidateDefinition(d, resolver)
}

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

import "arcoris.dev/apimachinery/api/identity"

// Resolver resolves API resource definitions by type and resource identities.
//
// Resolver belongs to package resource because Definition belongs to package
// resource. Concrete catalogs belong to higher composition packages. This lets
// API layers depend on resource contracts without depending on a concrete
// mutable catalog implementation.
type Resolver interface {
	// ResolveResource returns the definition registered for a group/resource key.
	ResolveResource(identity.GroupResource) (Definition, bool)

	// ResolveKind returns the definition registered for a group/kind key.
	ResolveKind(identity.GroupKind) (Definition, bool)

	// ResolveVersionResource returns the definition and exact version registered
	// for a group/version/resource key.
	ResolveVersionResource(identity.GroupVersionResource) (Definition, VersionDefinition, bool)

	// ResolveVersionKind returns the definition and exact version registered for
	// a group/version/kind key.
	ResolveVersionKind(identity.GroupVersionKind) (Definition, VersionDefinition, bool)
}

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
	"sync"

	"arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
)

// Catalog owns API resource definition descriptors.
//
// Catalog is zero-value usable, concurrency-safe, and explicitly owner-created.
// It deliberately has no package-level singleton or init-time registration
// path. That keeps API descriptor composition separate from runtime schemes,
// codecs, REST routing, storage, watches, controllers, and provider lifecycle
// behavior.
//
// Catalog must not be copied after first use.
type Catalog struct {
	// mu protects all catalog storage and indexes.
	//
	// Readers may resolve concurrently. Writers are serialized because
	// RegisterMany validates against a stable candidate catalog before
	// committing changes.
	mu sync.RWMutex

	// resolver validates resource surface TypeRef roots during registration.
	//
	// The field is intentionally set only by New or by the zero value. Changing
	// it after registration would weaken catalog invariants.
	resolver types.Resolver

	// defsByResource is the primary source of truth for registered definitions.
	defsByResource map[identity.GroupResource]resource.Definition

	// order preserves stable registration order for enumeration.
	order []identity.GroupResource

	// resourceByKind indexes version-independent kind identities.
	resourceByKind map[identity.GroupKind]identity.GroupResource

	// versionByResource indexes concrete version/resource identities.
	versionByResource map[identity.GroupVersionResource]versionRef

	// versionByKind indexes concrete version/kind identities.
	versionByKind map[identity.GroupVersionKind]versionRef
}

// versionRef points from a concrete version identity back to the primary
// resource-family record.
type versionRef struct {
	// resource identifies the primary Definition in defsByResource.
	resource identity.GroupResource

	// version identifies the exact VersionDefinition within the Definition.
	version identity.Version
}

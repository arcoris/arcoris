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

package meta

import (
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/meta/finalizer"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/meta/owner"
	"arcoris.dev/apimachinery/api/meta/stamp"
)

// ObjectMeta represents the metadata field of an API object.
//
// It contains identity metadata, system stamps, lifecycle metadata, labels,
// annotations, owner references, and finalizers. It does not contain desired or
// observed payloads, resource definitions, selectors, conditions, storage keys,
// or managed fields.
type ObjectMeta struct {
	// Name is the concrete metadata name of the object.
	Name metaidentity.Name
	// GenerateName is a server-side generation prefix, not a generated name.
	GenerateName metaidentity.NamePrefix
	// Namespace is absent when empty; this package does not default it.
	Namespace metaidentity.Namespace
	// UID identifies one concrete object incarnation.
	UID metaidentity.UID

	// ResourceVersion is an opaque concurrency/change token.
	ResourceVersion stamp.ResourceVersion
	// Generation tracks desired-state changes for one object.
	Generation stamp.Generation

	// CreatedAt records metadata creation time when a higher layer assigned it.
	CreatedAt stamp.Timestamp
	// Deletion records a deletion request without executing deletion.
	Deletion *stamp.Deletion

	// Labels are classification metadata for future indexing/selection layers.
	Labels labels.Set
	// Annotations are non-indexed metadata for tools and integrations.
	Annotations annotations.Set

	// OwnerReferences records ownership metadata without resolving owners.
	OwnerReferences owner.List
	// Finalizers names cleanup blockers; this package never runs cleanup.
	Finalizers finalizer.Set
}

// IsZero reports whether all object metadata fields are absent.
func (m ObjectMeta) IsZero() bool {
	return m.Name.IsZero() &&
		m.GenerateName.IsZero() &&
		m.Namespace.IsZero() &&
		m.UID.IsZero() &&
		m.ResourceVersion.IsZero() &&
		m.Generation.IsZero() &&
		m.CreatedAt.IsZero() &&
		(m.Deletion == nil || m.Deletion.IsZero()) &&
		m.Labels.IsZero() &&
		m.Annotations.IsZero() &&
		m.OwnerReferences.IsZero() &&
		m.Finalizers.IsZero()
}

// ObjectName returns the namespace/name portion of object metadata.
func (m ObjectMeta) ObjectName() metaidentity.ObjectName {
	return metaidentity.ObjectName{Namespace: m.Namespace, Name: m.Name}
}

// ObjectIdentity returns the namespace/name/UID portion of object metadata.
func (m ObjectMeta) ObjectIdentity() metaidentity.ObjectIdentity {
	return metaidentity.ObjectIdentity{Namespace: m.Namespace, Name: m.Name, UID: m.UID}
}

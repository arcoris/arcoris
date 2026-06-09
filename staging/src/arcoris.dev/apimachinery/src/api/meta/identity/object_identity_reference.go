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

package identity

import apiidentity "arcoris.dev/apimachinery/api/identity"

// ObjectIdentityReference describes another API object by API identity and UID.
//
// ObjectIdentityReference is UID-pinned: it identifies one concrete object
// incarnation rather than any future object that may reuse the same
// namespace/name. Owner metadata should use this type when ownership must be
// tied to a concrete owner object. Lifecycle policy, controller ownership, and
// block-owner-deletion behavior belong to api/meta/owner, not this raw identity
// package.
type ObjectIdentityReference struct {
	// APIVersion identifies the referenced object's API group/version.
	APIVersion apiidentity.GroupVersion `json:"apiVersion"`
	// Kind identifies the referenced object's API kind.
	Kind apiidentity.Kind `json:"kind"`
	// Namespace is optional reference metadata and is not defaulted here.
	Namespace Namespace `json:"namespace,omitempty"`
	// Name identifies the referenced object within its resource scope.
	Name Name `json:"name"`
	// UID pins the reference to one concrete object incarnation.
	UID UID `json:"uid"`
}

// IsZero reports whether all reference fields are absent.
func (r ObjectIdentityReference) IsZero() bool {
	return r.APIVersion.IsZero() &&
		r.Kind.IsZero() &&
		r.Namespace.IsZero() &&
		r.Name.IsZero() &&
		r.UID.IsZero()
}

// ObjectReference returns the name-only typed reference portion.
func (r ObjectIdentityReference) ObjectReference() ObjectReference {
	return ObjectReference{
		APIVersion: r.APIVersion,
		Kind:       r.Kind,
		Namespace:  r.Namespace,
		Name:       r.Name,
	}
}

// ObjectName returns the namespace/name portion of the reference.
func (r ObjectIdentityReference) ObjectName() ObjectName {
	return ObjectName{Namespace: r.Namespace, Name: r.Name}
}

// ObjectIdentity returns the namespace/name/UID portion of the reference.
func (r ObjectIdentityReference) ObjectIdentity() ObjectIdentity {
	return ObjectIdentity{Namespace: r.Namespace, Name: r.Name, UID: r.UID}
}

// GroupVersionKind returns the API group/version/kind portion of the reference.
func (r ObjectIdentityReference) GroupVersionKind() apiidentity.GroupVersionKind {
	return r.APIVersion.WithKind(r.Kind)
}

// String returns diagnostic text for the UID-pinned typed object reference.
//
// The result is intentionally diagnostic only. Storage keys, route keys, cache
// keys, and watch topics need explicit formats in higher layers because they
// also depend on resource descriptors, scope, and backend policy.
func (r ObjectIdentityReference) String() string {
	return r.GroupVersionKind().String() + " " + r.ObjectIdentity().String()
}

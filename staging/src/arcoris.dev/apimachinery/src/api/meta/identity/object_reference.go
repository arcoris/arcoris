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

// ObjectReference describes another API object by API identity and metadata
// identity.
//
// It does not resolve the object, check existence, or encode a route/storage
// key. It is reference metadata only.
type ObjectReference struct {
	// APIVersion identifies the referenced object's API group/version.
	APIVersion apiidentity.GroupVersion
	// Kind identifies the referenced object's API kind.
	Kind apiidentity.Kind
	// Namespace is optional reference metadata and is not defaulted here.
	Namespace Namespace
	// Name identifies the referenced object within its resource scope.
	Name Name
	// UID pins the reference to one concrete object incarnation.
	UID UID
}

// IsZero reports whether all reference fields are absent.
func (r ObjectReference) IsZero() bool {
	return r.APIVersion.IsZero() &&
		r.Kind.IsZero() &&
		r.Namespace.IsZero() &&
		r.Name.IsZero() &&
		r.UID.IsZero()
}

// ObjectName returns the namespace/name portion of the reference.
func (r ObjectReference) ObjectName() ObjectName {
	return ObjectName{Namespace: r.Namespace, Name: r.Name}
}

// ObjectIdentity returns the namespace/name/UID portion of the reference.
func (r ObjectReference) ObjectIdentity() ObjectIdentity {
	return ObjectIdentity{Namespace: r.Namespace, Name: r.Name, UID: r.UID}
}

// GroupVersionKind returns the API group/version/kind portion of the reference.
func (r ObjectReference) GroupVersionKind() apiidentity.GroupVersionKind {
	return r.APIVersion.WithKind(r.Kind)
}

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

import apiidentity "arcoris.dev/apimachinery/api/identity"

// TypeMeta represents object representation fields apiVersion and kind.
//
// It bridges object metadata representation to api/identity.GroupVersionKind.
// It does not contain resource, plural, scope, subresource, storage, or routing
// information.
type TypeMeta struct {
	// APIVersion is the canonical API group/version of the represented object.
	APIVersion apiidentity.GroupVersion
	// Kind is the API kind name of the represented object.
	Kind apiidentity.Kind
}

// FromGroupVersionKind constructs TypeMeta from a full API kind identity.
func FromGroupVersionKind(gvk apiidentity.GroupVersionKind) TypeMeta {
	return TypeMeta{APIVersion: gvk.GroupVersion(), Kind: gvk.Kind}
}

// GroupVersionKind returns the full API kind identity represented by TypeMeta.
func (m TypeMeta) GroupVersionKind() apiidentity.GroupVersionKind {
	return m.APIVersion.WithKind(m.Kind)
}

// IsZero reports whether apiVersion and kind are both absent.
func (m TypeMeta) IsZero() bool {
	return m.APIVersion.IsZero() && m.Kind.IsZero()
}

// String returns canonical diagnostic text for the represented kind identity.
func (m TypeMeta) String() string {
	return m.GroupVersionKind().String()
}

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

package objectstore

import metaidentity "arcoris.dev/apimachinery/api/meta/identity"

// ListScope defines structural filtering for objectstore.List.
//
// The store does not decide whether a resource is globally scoped or
// namespaced. It only applies this structural filter to object keys.
type ListScope struct {
	// kind is the requested structural list scope.
	kind ListScopeKind

	// namespace is present only when kind is ListScopeNamespace.
	namespace metaidentity.Namespace
}

// AllNamespaces returns a scope that matches every namespace for the resource.
func AllNamespaces() ListScope {
	return ListScope{kind: ListScopeAll}
}

// InNamespace validates namespace and returns a namespace-limited scope.
func InNamespace(namespace metaidentity.Namespace) (ListScope, error) {
	scope := ListScope{kind: ListScopeNamespace, namespace: namespace}
	if err := validateListScope(scope); err != nil {
		return ListScope{}, err
	}

	return scope, nil
}

// MustNamespace returns a namespace-limited scope or panics when namespace is invalid.
func MustNamespace(namespace metaidentity.Namespace) ListScope {
	scope, err := InNamespace(namespace)
	if err != nil {
		panic(err)
	}

	return scope
}

// Kind returns the structural scope kind.
func (s ListScope) Kind() ListScopeKind {
	return s.kind
}

// Namespace returns the namespace carried by a namespace-limited scope.
func (s ListScope) Namespace() metaidentity.Namespace {
	return s.namespace
}

// IsZero reports whether s has no scope kind.
func (s ListScope) IsZero() bool {
	return s.kind == 0 && s.namespace.IsZero()
}

// IsAllNamespaces reports whether s matches all namespaces.
func (s ListScope) IsAllNamespaces() bool {
	return s.kind == ListScopeAll
}

// IsNamespace reports whether s matches one namespace.
func (s ListScope) IsNamespace() bool {
	return s.kind == ListScopeNamespace
}

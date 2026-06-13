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

// ListScopeKind identifies structural collection-list scope.
type ListScopeKind uint8

const (
	// ListScopeAll lists all object names for the resource, regardless of namespace.
	ListScopeAll ListScopeKind = iota + 1

	// ListScopeNamespace lists only object names in one namespace.
	ListScopeNamespace
)

// IsValid reports whether k is one of the supported structural scopes.
func (k ListScopeKind) IsValid() bool {
	return k == ListScopeAll || k == ListScopeNamespace
}

// String returns stable diagnostic text for k.
func (k ListScopeKind) String() string {
	switch k {
	case ListScopeAll:
		return "all"
	case ListScopeNamespace:
		return "namespace"
	default:
		return "unknown"
	}
}

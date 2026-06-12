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

// Scope identifies the future instance-addressing scope of a resource family.
//
// Scope is part of the resource definition contract. It does not define object
// metadata, namespace value types, authorization, tenancy, storage partitioning,
// or routing behavior.
type Scope uint8

// Supported Scope values.
const (
	// ScopeInvalid is the zero value and is never a valid resource scope.
	ScopeInvalid Scope = iota
	// ScopeGlobal means future instances live in one global identity space.
	ScopeGlobal
	// ScopeNamespaced means future instances are partitioned by namespace.
	ScopeNamespaced
)

// IsZero reports whether s is the invalid zero scope.
func (s Scope) IsZero() bool { return s == ScopeInvalid }

// IsValid reports whether s is a supported scope value.
func (s Scope) IsValid() bool { return s == ScopeGlobal || s == ScopeNamespaced }

// String returns the canonical diagnostic and encoding text for s.
func (s Scope) String() string {
	switch s {
	case ScopeGlobal:
		return scopeTextGlobal
	case ScopeNamespaced:
		return scopeTextNamespaced
	case ScopeInvalid:
		return scopeTextInvalid
	default:
		return scopeTextUnknown
	}
}

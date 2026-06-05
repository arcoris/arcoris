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

package admission

// maxComponentKindLength keeps kind values short enough for stable diagnostic
// surfaces while leaving room for open-world domain-specific names.
const maxComponentKindLength = 64

// ComponentKind is an open-world coarse class for admission components.
//
// Kind groups a component by role, while ComponentID identifies the component
// more precisely. The value is stable metadata and must not contain runtime
// instance data.
type ComponentKind string

// MustComponentKind returns value as a ComponentKind or panics when value is
// invalid.
//
// The helper is intended for package-level constants and tests where an invalid
// literal is a programming error.
func MustComponentKind(value string) ComponentKind {
	kind := ComponentKind(value)
	if !kind.IsValid() {
		panic("admission.ComponentKind: invalid component kind")
	}
	return kind
}

// IsValid reports whether k is a valid lower_snake_case component kind.
//
// Custom kinds are allowed, but they must remain short, stable, ASCII
// identifiers suitable for logs, docs, and metrics dimensions owned elsewhere.
func (k ComponentKind) IsValid() bool {
	return validLowerSnakeIdentifier(string(k), maxComponentKindLength)
}

// String returns k as a string.
//
// The method intentionally performs no validation. Registry sorting and
// diagnostics can call it on both valid and invalid values without changing
// control flow.
func (k ComponentKind) String() string {
	return string(k)
}

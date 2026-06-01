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

package fieldpath

// ElementKind identifies one semantic step in a structured field path.
type ElementKind uint8

const (
	// ElementInvalid is the zero element kind and is never valid.
	ElementInvalid ElementKind = iota
	// ElementField addresses one fixed object field.
	ElementField
	// ElementKey addresses one dynamic map key.
	ElementKey
	// ElementIndex addresses one physical list position.
	ElementIndex
	// ElementSelector addresses one associative-list element identity.
	ElementSelector
)

// IsValid reports whether k identifies a supported element kind.
func (k ElementKind) IsValid() bool {
	return k >= ElementField && k <= ElementSelector
}

// String returns a stable diagnostic name for k.
func (k ElementKind) String() string {
	switch k {
	case ElementInvalid:
		return "invalid"
	case ElementField:
		return "field"
	case ElementKey:
		return "key"
	case ElementIndex:
		return "index"
	case ElementSelector:
		return "selector"
	default:
		return "unknown"
	}
}

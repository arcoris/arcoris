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

package codecselection

import "slices"

// Parameters is a deterministic immutable set of content-type parameters.
//
// Parameters are sorted by normalized name. Duplicate names are rejected because
// they would produce ambiguous selection keys.
type Parameters struct {
	// items stores normalized parameters sorted by name.
	items []Parameter
}

// NewParameters validates, normalizes, and sorts parameters by name.
func NewParameters(items ...Parameter) (Parameters, error) {
	return normalizeParametersAt(pathParameters, items)
}

// MustParameters returns normalized Parameters or panics when input is invalid.
func MustParameters(items ...Parameter) Parameters {
	parameters, err := NewParameters(items...)
	if err != nil {
		panic(err)
	}

	return parameters
}

// IsZero reports whether p contains no parameters.
func (p Parameters) IsZero() bool {
	return len(p.items) == 0
}

// Len returns the parameter count.
func (p Parameters) Len() int {
	return len(p.items)
}

// Items returns detached normalized parameters in deterministic order.
func (p Parameters) Items() []Parameter {
	return slices.Clone(p.items)
}

// Equal reports whether p and other contain the same normalized parameters.
func (p Parameters) Equal(other Parameters) bool {
	return p.key() == other.key()
}

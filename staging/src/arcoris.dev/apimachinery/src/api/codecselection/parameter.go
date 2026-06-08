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

// Parameter is one normalized content-type key parameter.
//
// Parameter names are lowercase ASCII tokens. Values are case-sensitive ASCII
// tokens. The type represents already-normalized selection key material, not raw
// MIME syntax.
type Parameter struct {
	// name is the normalized lowercase parameter name.
	name string

	// value is the normalized parameter value.
	value string
}

// NewParameter validates and normalizes one content-type parameter.
func NewParameter(name string, value string) (Parameter, error) {
	return normalizeParameterAt(pathParameter, Parameter{name: name, value: value})
}

// MustParameter returns a normalized Parameter or panics when input is invalid.
func MustParameter(name string, value string) Parameter {
	parameter, err := NewParameter(name, value)
	if err != nil {
		panic(err)
	}

	return parameter
}

// IsZero reports whether p contains no name and no value.
func (p Parameter) IsZero() bool {
	return p.name == "" && p.value == ""
}

// Name returns the normalized parameter name.
func (p Parameter) Name() string {
	return p.name
}

// Value returns the normalized parameter value.
func (p Parameter) Value() string {
	return p.value
}

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

// MarshalText returns the canonical group text after validation.
//
// Invalid direct literals are rejected so encoders cannot silently publish
// non-canonical group identities.
func (g Group) MarshalText() ([]byte, error) {
	return marshalText(g.String(), g.Validate)
}

// UnmarshalText parses a canonical group text value.
//
// The method uses ParseGroup exactly, so it has the same strict no-trim and
// no-normalization behavior as the parser.
func (g *Group) UnmarshalText(data []byte) error {
	if g == nil {
		return nilReceiver(identityNameGroup)
	}

	parsed, err := ParseGroup(string(data))
	if err != nil {
		return err
	}

	*g = parsed
	return nil
}

// MarshalJSON returns the canonical group as a JSON string.
//
// Identity JSON encoding is scalar-only. Object forms are reserved for higher
// layers that define their own resource or metadata documents.
func (g Group) MarshalJSON() ([]byte, error) {
	return marshalJSONString(g.String(), g.Validate)
}

// UnmarshalJSON parses a canonical group from a JSON string.
//
// Null, objects, arrays, numbers, and booleans are rejected instead of being
// interpreted as the core group.
func (g *Group) UnmarshalJSON(data []byte) error {
	if g == nil {
		return nilReceiver(identityNameGroup)
	}

	value, err := unmarshalJSONString(identityNameGroup, data)
	if err != nil {
		return err
	}

	return g.UnmarshalText([]byte(value))
}

/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package schema

// Group identifies an API group.
//
// The empty group is the canonical spelling of the core API group. Non-empty
// groups must be DNS-1123 subdomains and are never case-normalized. Runtime
// and discovery layers should treat this value as an exact wire token, not as a
// display name or user input that can be normalized later.
type Group string

// ParseGroup parses an API group in canonical form.
//
// An empty string is accepted and represents the core group. Non-empty values
// must be DNS-1123 subdomains. The parser does not trim whitespace, does not
// lower-case input, and does not accept path-like or underscore-separated
// alternatives.
func ParseGroup(val string) (Group, error) {
	group := Group(val)
	if err := group.Validate(); err != nil {
		return "", err
	}
	return group, nil
}

// String returns the canonical group string without revalidating it.
//
// Call Validate at trust boundaries when the value may have been constructed
// directly instead of returned by ParseGroup.
func (g Group) String() string {
	return string(g)
}

// Validate checks that the group is empty or a DNS-1123 subdomain.
//
// Empty is accepted only because it is the explicit core-group identity. Every
// non-empty label must be lowercase ASCII DNS-1123.
func (g Group) Validate() error {
	return validateGroupValue(string(g))
}

// IsZero reports whether the group is empty.
//
// For Group, zero and the core group are the same canonical value. Composite
// identities still require their own non-group fields to be present.
func (g Group) IsZero() bool {
	return g == ""
}

// MarshalText returns the canonical group string after validation.
//
// Invalid direct literals are rejected here so text encoding cannot silently
// publish non-canonical identifiers.
func (g Group) MarshalText() ([]byte, error) {
	if err := g.Validate(); err != nil {
		return nil, err
	}
	return []byte(g.String()), nil
}

// UnmarshalText parses a canonical group string.
//
// The method uses ParseGroup exactly, so it inherits the same strict no-trim
// behavior as the public parser.
func (g *Group) UnmarshalText(data []byte) error {
	if g == nil {
		return nilUnmarshalReceiver("group")
	}
	parsed, err := ParseGroup(string(data))
	if err != nil {
		return err
	}
	*g = parsed
	return nil
}

// MarshalJSON returns the canonical group string as a JSON scalar.
//
// The JSON representation of every schema identifier is a scalar string; there
// is no object form for atomic or composite identifiers.
func (g Group) MarshalJSON() ([]byte, error) {
	return marshalJSONString("group", g.String(), g.Validate)
}

// UnmarshalJSON parses a canonical group string from a JSON scalar.
//
// Non-string JSON, including null, is rejected instead of being interpreted as
// the zero or core group.
func (g *Group) UnmarshalJSON(data []byte) error {
	if g == nil {
		return nilUnmarshalReceiver("group")
	}
	val, err := unmarshalJSONString("group", data)
	if err != nil {
		return err
	}
	return g.UnmarshalText([]byte(val))
}

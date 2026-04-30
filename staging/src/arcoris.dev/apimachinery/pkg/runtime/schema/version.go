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

// Version identifies an API version.
//
// ARCORIS accepts only strict version tokens: vN, vNalphaM, or vNbetaM. Numeric
// components are base-10 ASCII digits with no leading zero except the literal
// version v0. Version values are protocol identifiers and must not be
// case-normalized or treated as semantic versions.
type Version string

// ParseVersion parses a strict ARCORIS API version token.
//
// The parser rejects empty input, uppercase variants, release-candidate tags,
// preview names, and numeric forms such as "1" or "v01". It does not trim
// whitespace.
func ParseVersion(value string) (Version, error) {
	version := Version(value)
	if err := version.Validate(); err != nil {
		return "", err
	}
	return version, nil
}

// String returns the canonical version string without revalidating it.
//
// Values returned by ParseVersion are canonical by construction. Direct
// literals should be checked with Validate before leaving a trust boundary.
func (v Version) String() string {
	return string(v)
}

// Validate checks that the version matches vN, vNalphaM, or vNbetaM.
//
// The method is the strict contract boundary for values constructed directly
// from strings.
func (v Version) Validate() error {
	return validateVersionValue(string(v))
}

// IsZero reports whether the version is empty.
//
// A zero Version is allowed as an optional-field sentinel, but complete
// identities such as GroupVersion reject a missing version.
func (v Version) IsZero() bool {
	return v == ""
}

// MarshalText returns the canonical version string after validation.
//
// Invalid direct literals are rejected instead of being serialized.
func (v Version) MarshalText() ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return []byte(v.String()), nil
}

// UnmarshalText parses a canonical version string.
//
// The method delegates to ParseVersion and therefore rejects all non-canonical
// spellings accepted by neither ARCORIS nor the version grammar.
func (v *Version) UnmarshalText(data []byte) error {
	if v == nil {
		return nilUnmarshalReceiver("version")
	}
	parsed, err := ParseVersion(string(data))
	if err != nil {
		return err
	}
	*v = parsed
	return nil
}

// MarshalJSON returns the canonical version string as a JSON scalar.
//
// JSON object forms are never emitted by schema identifiers.
func (v Version) MarshalJSON() ([]byte, error) {
	return marshalJSONString("version", v.String(), v.Validate)
}

// UnmarshalJSON parses a canonical version string from a JSON scalar.
//
// Non-string JSON, including null, is rejected before version parsing.
func (v *Version) UnmarshalJSON(data []byte) error {
	if v == nil {
		return nilUnmarshalReceiver("version")
	}
	value, err := unmarshalJSONString("version", data)
	if err != nil {
		return err
	}
	return v.UnmarshalText([]byte(value))
}

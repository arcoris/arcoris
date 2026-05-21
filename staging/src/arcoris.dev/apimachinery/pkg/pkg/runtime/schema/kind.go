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

// Kind identifies an API object kind.
//
// Kinds are PascalCase-like ASCII identifiers. They are never normalized and
// must start with an uppercase ASCII letter. A Kind names a type schema; it is
// not a resource collection name and does not contain group or version
// segments.
type Kind string

// ParseKind parses a canonical API kind.
//
// The parser accepts only ASCII letters and digits after the initial uppercase
// letter. It rejects separators, whitespace, Unicode display names, and
// lowercase starts.
func ParseKind(val string) (Kind, error) {
	kind := Kind(val)
	if err := kind.Validate(); err != nil {
		return "", err
	}
	return kind, nil
}

// String returns the canonical kind string without revalidating it.
//
// Use Validate when accepting a direct Kind literal from another package.
func (k Kind) String() string {
	return string(k)
}

// Validate checks that the kind is a PascalCase-like ASCII identifier.
//
// This is the boundary that prevents path segments, dotted identities, and
// localized names from entering schema type identity.
func (k Kind) Validate() error {
	return validateKindValue(string(k))
}

// IsZero reports whether the kind is empty.
//
// Zero can be useful while constructing optional values, but complete type
// identities reject a missing kind.
func (k Kind) IsZero() bool {
	return k == ""
}

// MarshalText returns the canonical kind string after validation.
//
// Invalid direct literals are rejected instead of being serialized.
func (k Kind) MarshalText() ([]byte, error) {
	if err := k.Validate(); err != nil {
		return nil, err
	}
	return []byte(k.String()), nil
}

// UnmarshalText parses a canonical kind string.
//
// The method uses the same strict parser as normal API-boundary code.
func (k *Kind) UnmarshalText(data []byte) error {
	if k == nil {
		return nilUnmarshalReceiver("kind")
	}
	parsed, err := ParseKind(string(data))
	if err != nil {
		return err
	}
	*k = parsed
	return nil
}

// MarshalJSON returns the canonical kind string as a JSON scalar.
//
// Schema identifiers serialize as strings, never as structured JSON objects.
func (k Kind) MarshalJSON() ([]byte, error) {
	return marshalJSONString("kind", k.String(), k.Validate)
}

// UnmarshalJSON parses a canonical kind string from a JSON scalar.
//
// Non-string JSON, including null, is rejected.
func (k *Kind) UnmarshalJSON(data []byte) error {
	if k == nil {
		return nilUnmarshalReceiver("kind")
	}
	val, err := unmarshalJSONString("kind", data)
	if err != nil {
		return err
	}
	return k.UnmarshalText([]byte(val))
}

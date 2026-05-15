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

// Subresource identifies an API subresource.
//
// The empty subresource is valid and means no subresource. Non-empty
// subresources use the same DNS-1123 single-label contract as Resource. A
// Subresource is always separate from Resource so callers cannot smuggle
// "pods/status" into a resource collection name.
type Subresource string

// ParseSubresource parses a canonical API subresource name.
//
// Empty input is accepted as the canonical absent-subresource value. Non-empty
// input must be a DNS-1123 single label and is not trimmed or normalized.
func ParseSubresource(val string) (Subresource, error) {
	subresource := Subresource(val)
	if err := subresource.Validate(); err != nil {
		return "", err
	}
	return subresource, nil
}

// String returns the canonical subresource string without revalidating it.
//
// Direct literals should be validated before crossing an API boundary.
func (s Subresource) String() string {
	return string(s)
}

// Validate checks that the subresource is empty or a DNS-1123 single label.
//
// Empty is valid only because it explicitly means "no subresource"; malformed
// non-empty values are rejected.
func (s Subresource) Validate() error {
	return validateDNS1123SingleLabel("subresource", string(s), true)
}

// IsZero reports whether the subresource is empty.
//
// For Subresource, zero and absence are the same canonical value.
func (s Subresource) IsZero() bool {
	return s == ""
}

// MarshalText returns the canonical subresource string after validation.
//
// The empty subresource marshals as an empty string because it is a valid
// absent-subresource identity.
func (s Subresource) MarshalText() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return []byte(s.String()), nil
}

// UnmarshalText parses a canonical subresource string.
//
// The method uses ParseSubresource and preserves the same strictness.
func (s *Subresource) UnmarshalText(data []byte) error {
	if s == nil {
		return nilUnmarshalReceiver("subresource")
	}
	parsed, err := ParseSubresource(string(data))
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// MarshalJSON returns the canonical subresource string as a JSON scalar.
//
// Empty subresources marshal as "" and non-string JSON is never emitted.
func (s Subresource) MarshalJSON() ([]byte, error) {
	return marshalJSONString("subresource", s.String(), s.Validate)
}

// UnmarshalJSON parses a canonical subresource string from a JSON scalar.
//
// Non-string JSON, including null, is rejected.
func (s *Subresource) UnmarshalJSON(data []byte) error {
	if s == nil {
		return nilUnmarshalReceiver("subresource")
	}
	val, err := unmarshalJSONString("subresource", data)
	if err != nil {
		return err
	}
	return s.UnmarshalText([]byte(val))
}

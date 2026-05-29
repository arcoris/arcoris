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

// MarshalText returns the canonical scope text.
//
// Scope scalar encoding mirrors api/identity scalar encoding. It is not a
// manifest format for Definition and does not imply JSON tags, DTOs, or
// resource-definition marshaling in this package.
func (s Scope) MarshalText() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return []byte(s.String()), nil
}

// UnmarshalText parses canonical scope text.
func (s *Scope) UnmarshalText(text []byte) error {
	if s == nil {
		return nilReceiver(pathScope)
	}
	parsed, err := ParseScope(string(text))
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

// MarshalJSON returns the canonical scope as a JSON string.
//
// The JSON shape is scalar-only. Definition intentionally has no JSON encoding
// here because manifest/codecs/exporters belong to higher API layers.
func (s Scope) MarshalJSON() ([]byte, error) {
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return marshalJSONString(s.String())
}

// UnmarshalJSON parses a JSON string scope.
func (s *Scope) UnmarshalJSON(data []byte) error {
	if s == nil {
		return nilReceiver(pathScope)
	}
	value, err := unmarshalJSONString(pathScope, data)
	if err != nil {
		return err
	}
	parsed, err := ParseScope(value)
	if err != nil {
		return err
	}
	*s = parsed
	return nil
}

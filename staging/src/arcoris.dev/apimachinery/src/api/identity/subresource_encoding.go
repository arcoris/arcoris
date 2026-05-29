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

// MarshalText returns the canonical subresource text after validation.
func (s Subresource) MarshalText() ([]byte, error) {
	return marshalText(s.String(), s.Validate)
}

// UnmarshalText parses a canonical subresource text value.
func (s *Subresource) UnmarshalText(data []byte) error {
	if s == nil {
		return nilReceiver(identityNameSubresource)
	}

	parsed, err := ParseSubresource(string(data))
	if err != nil {
		return err
	}

	*s = parsed
	return nil
}

// MarshalJSON returns the canonical subresource as a JSON string.
func (s Subresource) MarshalJSON() ([]byte, error) {
	return marshalJSONString(s.String(), s.Validate)
}

// UnmarshalJSON parses a canonical subresource from a JSON string.
func (s *Subresource) UnmarshalJSON(data []byte) error {
	if s == nil {
		return nilReceiver(identityNameSubresource)
	}

	value, err := unmarshalJSONString(identityNameSubresource, data)
	if err != nil {
		return err
	}

	return s.UnmarshalText([]byte(value))
}

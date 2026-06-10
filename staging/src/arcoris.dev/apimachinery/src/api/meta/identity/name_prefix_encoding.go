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

// MarshalText validates and encodes the name prefix as text.
func (p NamePrefix) MarshalText() ([]byte, error) {
	return marshalText(p.String(), p.ValidateLexical)
}

// UnmarshalText decodes and validates a text name prefix.
func (p *NamePrefix) UnmarshalText(data []byte) error {
	if p == nil {
		return nilReceiver("namePrefix")
	}

	value, err := ParseNamePrefix(string(data))
	if err != nil {
		return err
	}

	*p = value
	return nil
}

// MarshalJSON validates and encodes the name prefix as one JSON string.
func (p NamePrefix) MarshalJSON() ([]byte, error) {
	return marshalJSONString(p.String(), p.ValidateLexical)
}

// UnmarshalJSON decodes and validates a JSON string name prefix.
func (p *NamePrefix) UnmarshalJSON(data []byte) error {
	if p == nil {
		return nilReceiver("namePrefix")
	}

	value, err := unmarshalJSONString("namePrefix", data)
	if err != nil {
		return err
	}

	parsed, err := ParseNamePrefix(value)
	if err != nil {
		return err
	}

	*p = parsed
	return nil
}

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

package finalizer

// MarshalText validates and encodes the finalizer name as scalar text.
func (n Name) MarshalText() ([]byte, error) {
	return marshalText(n.String(), n.Validate)
}

// UnmarshalText decodes and validates scalar finalizer name text.
func (n *Name) UnmarshalText(data []byte) error {
	if n == nil {
		return nilReceiver("finalizer.name")
	}

	value, err := ParseName(string(data))
	if err != nil {
		return err
	}

	*n = value
	return nil
}

// MarshalJSON validates and encodes the finalizer name as one JSON string.
func (n Name) MarshalJSON() ([]byte, error) {
	return marshalJSONString(n.String(), n.Validate)
}

// UnmarshalJSON decodes one JSON string and rejects null or non-string input.
func (n *Name) UnmarshalJSON(data []byte) error {
	if n == nil {
		return nilReceiver("finalizer.name")
	}

	value, err := unmarshalJSONString("finalizer.name", data)
	if err != nil {
		return err
	}

	parsed, err := ParseName(value)
	if err != nil {
		return err
	}

	*n = parsed
	return nil
}

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

// MarshalText validates and encodes the namespace as text.
func (n Namespace) MarshalText() ([]byte, error) {
	return marshalText(n.String(), n.ValidateLexical)
}

// UnmarshalText decodes and validates a text namespace.
func (n *Namespace) UnmarshalText(data []byte) error {
	if n == nil {
		return nilReceiver("namespace")
	}

	value, err := ParseNamespace(string(data))
	if err != nil {
		return err
	}

	*n = value
	return nil
}

// MarshalJSON validates and encodes the namespace as one JSON string.
func (n Namespace) MarshalJSON() ([]byte, error) {
	return marshalJSONString(n.String(), n.ValidateLexical)
}

// UnmarshalJSON decodes and validates a JSON string namespace.
func (n *Namespace) UnmarshalJSON(data []byte) error {
	if n == nil {
		return nilReceiver("namespace")
	}

	value, err := unmarshalJSONString("namespace", data)
	if err != nil {
		return err
	}

	parsed, err := ParseNamespace(value)
	if err != nil {
		return err
	}

	*n = parsed
	return nil
}

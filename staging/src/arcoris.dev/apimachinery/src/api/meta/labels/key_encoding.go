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

package labels

// MarshalText validates and encodes the key as text.
func (k Key) MarshalText() ([]byte, error) {
	return marshalText(k.String(), k.Validate)
}

// UnmarshalText decodes and validates a text key.
func (k *Key) UnmarshalText(data []byte) error {
	if k == nil {
		return nilReceiver("label.key")
	}

	value, err := ParseKey(string(data))
	if err != nil {
		return err
	}

	*k = value
	return nil
}

// MarshalJSON validates and encodes the key as one JSON string.
func (k Key) MarshalJSON() ([]byte, error) {
	return marshalJSONString(k.String(), k.Validate)
}

// UnmarshalJSON decodes and validates a JSON string key.
func (k *Key) UnmarshalJSON(data []byte) error {
	if k == nil {
		return nilReceiver("label.key")
	}

	value, err := unmarshalJSONString("label.key", data)
	if err != nil {
		return err
	}

	parsed, err := ParseKey(value)
	if err != nil {
		return err
	}

	*k = parsed
	return nil
}

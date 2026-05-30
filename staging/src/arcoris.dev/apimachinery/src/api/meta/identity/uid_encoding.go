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

// MarshalText validates and encodes the UID as text.
func (u UID) MarshalText() ([]byte, error) {
	return marshalText(u.String(), u.Validate)
}

// UnmarshalText decodes and validates a text UID.
func (u *UID) UnmarshalText(data []byte) error {
	if u == nil {
		return nilReceiver("uid")
	}

	value, err := ParseUID(string(data))
	if err != nil {
		return err
	}

	*u = value
	return nil
}

// MarshalJSON validates and encodes the UID as one JSON string.
func (u UID) MarshalJSON() ([]byte, error) {
	return marshalJSONString(u.String(), u.Validate)
}

// UnmarshalJSON decodes and validates a JSON string UID.
func (u *UID) UnmarshalJSON(data []byte) error {
	if u == nil {
		return nilReceiver("uid")
	}

	value, err := unmarshalJSONString("uid", data)
	if err != nil {
		return err
	}

	parsed, err := ParseUID(value)
	if err != nil {
		return err
	}

	*u = parsed
	return nil
}

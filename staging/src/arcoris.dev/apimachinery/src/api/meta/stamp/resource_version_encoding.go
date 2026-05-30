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

package stamp

// MarshalText validates and encodes the resource version as text.
func (v ResourceVersion) MarshalText() ([]byte, error) {
	return marshalText(v.String(), v.Validate)
}

// UnmarshalText decodes and validates a text resource version.
func (v *ResourceVersion) UnmarshalText(data []byte) error {
	if v == nil {
		return nilReceiver("resourceVersion")
	}

	value, err := ParseResourceVersion(string(data))
	if err != nil {
		return err
	}

	*v = value
	return nil
}

// MarshalJSON validates and encodes the resource version as one JSON string.
func (v ResourceVersion) MarshalJSON() ([]byte, error) {
	return marshalJSONString(v.String(), v.Validate)
}

// UnmarshalJSON decodes and validates a JSON string resource version.
func (v *ResourceVersion) UnmarshalJSON(data []byte) error {
	if v == nil {
		return nilReceiver("resourceVersion")
	}

	value, err := unmarshalJSONString("resourceVersion", data)
	if err != nil {
		return err
	}

	parsed, err := ParseResourceVersion(value)
	if err != nil {
		return err
	}

	*v = parsed
	return nil
}

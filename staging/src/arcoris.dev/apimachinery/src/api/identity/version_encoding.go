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

// MarshalText returns the canonical version text after validation.
func (v Version) MarshalText() ([]byte, error) {
	return marshalText(v.String(), v.Validate)
}

// UnmarshalText parses a canonical version text value.
func (v *Version) UnmarshalText(data []byte) error {
	if v == nil {
		return nilReceiver(identityNameVersion)
	}

	parsed, err := ParseVersion(string(data))
	if err != nil {
		return err
	}

	*v = parsed
	return nil
}

// MarshalJSON returns the canonical version as a JSON string.
func (v Version) MarshalJSON() ([]byte, error) {
	return marshalJSONString(v.String(), v.Validate)
}

// UnmarshalJSON parses a canonical version from a JSON string.
func (v *Version) UnmarshalJSON(data []byte) error {
	if v == nil {
		return nilReceiver(identityNameVersion)
	}

	value, err := unmarshalJSONString(identityNameVersion, data)
	if err != nil {
		return err
	}

	return v.UnmarshalText([]byte(value))
}

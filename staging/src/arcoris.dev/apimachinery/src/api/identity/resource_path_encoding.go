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

// MarshalText returns the canonical resource path after validation.
func (rp ResourcePath) MarshalText() ([]byte, error) {
	return marshalText(rp.String(), rp.Validate)
}

// UnmarshalText parses a canonical resource path text value.
func (rp *ResourcePath) UnmarshalText(data []byte) error {
	if rp == nil {
		return nilReceiver(identityNameResourcePath)
	}

	parsed, err := ParseResourcePath(string(data))
	if err != nil {
		return err
	}

	*rp = parsed
	return nil
}

// MarshalJSON returns the canonical resource path as a JSON string.
func (rp ResourcePath) MarshalJSON() ([]byte, error) {
	return marshalJSONString(rp.String(), rp.Validate)
}

// UnmarshalJSON parses a canonical resource path from a JSON string.
func (rp *ResourcePath) UnmarshalJSON(data []byte) error {
	if rp == nil {
		return nilReceiver(identityNameResourcePath)
	}

	value, err := unmarshalJSONString(identityNameResourcePath, data)
	if err != nil {
		return err
	}

	return rp.UnmarshalText([]byte(value))
}

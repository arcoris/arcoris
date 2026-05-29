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

// MarshalText returns the canonical group/version text after validation.
func (gv GroupVersion) MarshalText() ([]byte, error) {
	return marshalText(gv.String(), gv.Validate)
}

// UnmarshalText parses a canonical group/version text value.
func (gv *GroupVersion) UnmarshalText(data []byte) error {
	if gv == nil {
		return nilReceiver(identityNameGroupVersion)
	}

	parsed, err := ParseGroupVersion(string(data))
	if err != nil {
		return err
	}

	*gv = parsed
	return nil
}

// MarshalJSON returns the canonical group/version identity as a JSON string.
func (gv GroupVersion) MarshalJSON() ([]byte, error) {
	return marshalJSONString(gv.String(), gv.Validate)
}

// UnmarshalJSON parses a canonical group/version identity from a JSON string.
func (gv *GroupVersion) UnmarshalJSON(data []byte) error {
	if gv == nil {
		return nilReceiver(identityNameGroupVersion)
	}

	value, err := unmarshalJSONString(identityNameGroupVersion, data)
	if err != nil {
		return err
	}

	return gv.UnmarshalText([]byte(value))
}

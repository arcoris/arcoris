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

// MarshalText returns the canonical resource text after validation.
func (r Resource) MarshalText() ([]byte, error) {
	return marshalText(r.String(), r.Validate)
}

// UnmarshalText parses a canonical resource text value.
func (r *Resource) UnmarshalText(data []byte) error {
	if r == nil {
		return nilReceiver(identityNameResource)
	}

	parsed, err := ParseResource(string(data))
	if err != nil {
		return err
	}

	*r = parsed
	return nil
}

// MarshalJSON returns the canonical resource as a JSON string.
func (r Resource) MarshalJSON() ([]byte, error) {
	return marshalJSONString(r.String(), r.Validate)
}

// UnmarshalJSON parses a canonical resource from a JSON string.
func (r *Resource) UnmarshalJSON(data []byte) error {
	if r == nil {
		return nilReceiver(identityNameResource)
	}

	value, err := unmarshalJSONString(identityNameResource, data)
	if err != nil {
		return err
	}

	return r.UnmarshalText([]byte(value))
}

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

// MarshalText returns the canonical group/version/kind text after validation.
func (gvk GroupVersionKind) MarshalText() ([]byte, error) {
	return marshalText(gvk.String(), gvk.Validate)
}

// UnmarshalText parses a canonical group/version/kind text value.
func (gvk *GroupVersionKind) UnmarshalText(data []byte) error {
	if gvk == nil {
		return nilReceiver(identityNameGroupVersionKind)
	}

	parsed, err := ParseGroupVersionKind(string(data))
	if err != nil {
		return err
	}

	*gvk = parsed
	return nil
}

// MarshalJSON returns the canonical group/version/kind identity as a JSON string.
func (gvk GroupVersionKind) MarshalJSON() ([]byte, error) {
	return marshalJSONString(gvk.String(), gvk.Validate)
}

// UnmarshalJSON parses a canonical group/version/kind identity from a JSON string.
func (gvk *GroupVersionKind) UnmarshalJSON(data []byte) error {
	if gvk == nil {
		return nilReceiver(identityNameGroupVersionKind)
	}

	value, err := unmarshalJSONString(identityNameGroupVersionKind, data)
	if err != nil {
		return err
	}

	return gvk.UnmarshalText([]byte(value))
}

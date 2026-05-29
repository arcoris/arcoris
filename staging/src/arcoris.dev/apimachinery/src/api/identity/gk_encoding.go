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

// MarshalText returns the canonical group/kind text after validation.
func (gk GroupKind) MarshalText() ([]byte, error) {
	return marshalText(gk.String(), gk.Validate)
}

// UnmarshalText parses a canonical group/kind text value.
func (gk *GroupKind) UnmarshalText(data []byte) error {
	if gk == nil {
		return nilReceiver(identityNameGroupKind)
	}

	parsed, err := ParseGroupKind(string(data))
	if err != nil {
		return err
	}

	*gk = parsed
	return nil
}

// MarshalJSON returns the canonical group/kind identity as a JSON string.
func (gk GroupKind) MarshalJSON() ([]byte, error) {
	return marshalJSONString(gk.String(), gk.Validate)
}

// UnmarshalJSON parses a canonical group/kind identity from a JSON string.
func (gk *GroupKind) UnmarshalJSON(data []byte) error {
	if gk == nil {
		return nilReceiver(identityNameGroupKind)
	}

	value, err := unmarshalJSONString(identityNameGroupKind, data)
	if err != nil {
		return err
	}

	return gk.UnmarshalText([]byte(value))
}

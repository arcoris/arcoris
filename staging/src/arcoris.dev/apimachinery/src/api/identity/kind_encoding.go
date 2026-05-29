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

// MarshalText returns the canonical kind text after validation.
func (k Kind) MarshalText() ([]byte, error) {
	return marshalText(k.String(), k.Validate)
}

// UnmarshalText parses a canonical kind text value.
func (k *Kind) UnmarshalText(data []byte) error {
	if k == nil {
		return nilReceiver(identityNameKind)
	}

	parsed, err := ParseKind(string(data))
	if err != nil {
		return err
	}

	*k = parsed
	return nil
}

// MarshalJSON returns the canonical kind as a JSON string.
func (k Kind) MarshalJSON() ([]byte, error) {
	return marshalJSONString(k.String(), k.Validate)
}

// UnmarshalJSON parses a canonical kind from a JSON string.
func (k *Kind) UnmarshalJSON(data []byte) error {
	if k == nil {
		return nilReceiver(identityNameKind)
	}

	value, err := unmarshalJSONString(identityNameKind, data)
	if err != nil {
		return err
	}

	return k.UnmarshalText([]byte(value))
}

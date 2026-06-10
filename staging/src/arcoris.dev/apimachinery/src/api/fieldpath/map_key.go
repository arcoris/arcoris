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

package fieldpath

// MapKey is one dynamic map key in a semantic path.
//
// Map keys are path tokens only. They must be non-empty, but descriptor key
// constraints, label grammar, and map value semantics belong to higher layers.
type MapKey string

// NewMapKey validates key as a dynamic map-key token.
func NewMapKey(key string) (MapKey, error) {
	mapKey := MapKey(key)
	if err := mapKey.ValidateStructure(); err != nil {
		return "", err
	}

	return mapKey, nil
}

// MustMapKey validates key or panics.
//
// It is intended for tests and static semantic-path declarations. Runtime
// callers that receive untrusted text should use NewMapKey.
func MustMapKey(key string) MapKey {
	mapKey, err := NewMapKey(key)
	if err != nil {
		panic(err)
	}

	return mapKey
}

// String returns the map key text.
func (k MapKey) String() string {
	return string(k)
}

// IsZero reports whether k is the absent map-key value.
func (k MapKey) IsZero() bool {
	return k == ""
}

// ValidateStructure checks the base fieldpath map-key invariant.
//
// It does not validate descriptor map-key constraints or resource-specific
// semantics.
func (k MapKey) ValidateStructure() error {
	if !k.IsZero() {
		return nil
	}

	return nested(
		ErrInvalidElement,
		ErrorReasonEmptyMapKey,
		"map key is empty",
		ErrEmptyMapKey,
	)
}

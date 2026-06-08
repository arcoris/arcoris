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

package objectstore

import "errors"

// IsValid reports whether k is structurally usable as an object store key.
func (k Key) IsValid() bool {
	return ValidateKey(k) == nil
}

// ValidateKey checks only the structural key shape required for storage.
func ValidateKey(key Key) error {
	if err := key.Resource.Validate(); err != nil {
		return errorFor(ReasonInvalidKey, key, 0, 0, errors.Join(ErrInvalidKey, err))
	}
	if err := key.Object.Validate(); err != nil {
		return errorFor(ReasonInvalidKey, key, 0, 0, errors.Join(ErrInvalidKey, err))
	}

	return nil
}

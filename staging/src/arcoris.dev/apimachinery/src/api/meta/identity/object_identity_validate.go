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

// Validate checks that an object identity has a valid object name and UID.
func (i ObjectIdentity) Validate() error {
	if err := i.ObjectName().Validate(); err != nil {
		return nested("objectIdentity.objectName", ErrInvalidObjectIdentity, err)
	}

	if i.UID.IsZero() {
		return invalid(
			"objectIdentity.uid",
			ErrInvalidObjectIdentity,
			ErrorReasonEmptyValue,
			"uid is required",
		)
	}

	if err := i.UID.Validate(); err != nil {
		return nested("objectIdentity.uid", ErrInvalidObjectIdentity, err)
	}

	return nil
}

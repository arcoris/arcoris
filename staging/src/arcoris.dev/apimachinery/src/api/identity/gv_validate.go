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

// Validate checks that the group/version identity is complete and canonical.
//
// Group may be the core group. Version is required because an API group without
// a version is not a complete API identity.
func (gv GroupVersion) Validate() error {
	if gv.Version.IsZero() {
		return invalid(
			identityNameGroupVersion,
			gv.String(),
			ErrorReasonEmptyValue,
			detailVersionRequired,
		)
	}

	if err := gv.Group.Validate(); err != nil {
		return invalidValue(identityNameGroupVersion, gv.String(), err)
	}

	if err := gv.Version.Validate(); err != nil {
		return invalidValue(identityNameGroupVersion, gv.String(), err)
	}

	return nil
}

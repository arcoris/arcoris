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

// Validate checks that the group/kind identity is complete and canonical.
//
// Group may be the core group. Kind is required because a group alone does not
// identify an API object kind.
func (gk GroupKind) Validate() error {
	if gk.Kind.IsZero() {
		return invalid(
			identityNameGroupKind,
			gk.String(),
			ErrorReasonEmptyValue,
			detailKindRequired,
		)
	}

	if err := gk.Group.Validate(); err != nil {
		return invalidValue(identityNameGroupKind, gk.String(), err)
	}

	if err := gk.Kind.Validate(); err != nil {
		return invalidValue(identityNameGroupKind, gk.String(), err)
	}

	return nil
}

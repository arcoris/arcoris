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

// Validate checks that the group/version/kind identity is complete and canonical.
func (gvk GroupVersionKind) Validate() error {
	if gvk.Version.IsZero() {
		return invalid(
			identityNameGroupVersionKind,
			gvk.String(),
			ErrorReasonEmptyValue,
			detailVersionRequired,
		)
	}

	if gvk.Kind.IsZero() {
		return invalid(
			identityNameGroupVersionKind,
			gvk.String(),
			ErrorReasonEmptyValue,
			detailKindRequired,
		)
	}

	if err := gvk.GroupVersion().Validate(); err != nil {
		return invalidValue(identityNameGroupVersionKind, gvk.String(), err)
	}

	if err := gvk.Kind.Validate(); err != nil {
		return invalidValue(identityNameGroupVersionKind, gvk.String(), err)
	}

	return nil
}

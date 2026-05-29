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

// Validate checks that the versioned resource path is complete and canonical.
func (gvrp GroupVersionResourcePath) Validate() error {
	if gvrp.Version.IsZero() {
		return invalid(
			identityNameGroupVersionResourcePath,
			gvrp.String(),
			ErrorReasonEmptyValue,
			detailVersionRequired,
		)
	}

	if gvrp.Resource.IsZero() {
		return invalid(
			identityNameGroupVersionResourcePath,
			gvrp.String(),
			ErrorReasonEmptyValue,
			detailResourceRequired,
		)
	}

	if err := gvrp.GroupVersion().Validate(); err != nil {
		return invalidValue(identityNameGroupVersionResourcePath, gvrp.String(), err)
	}

	if err := gvrp.ResourcePath().Validate(); err != nil {
		return invalidValue(identityNameGroupVersionResourcePath, gvrp.String(), err)
	}

	return nil
}

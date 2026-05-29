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

// Validate checks that the group/resource identity is complete and canonical.
//
// Group may be the core group. Resource is required because a group alone does
// not identify an API collection.
func (gr GroupResource) Validate() error {
	if gr.Resource.IsZero() {
		return invalid(
			identityNameGroupResource,
			gr.String(),
			ErrorReasonEmptyValue,
			detailResourceRequired,
		)
	}

	if err := gr.Group.Validate(); err != nil {
		return invalidValue(identityNameGroupResource, gr.String(), err)
	}

	if err := gr.Resource.Validate(); err != nil {
		return invalidValue(identityNameGroupResource, gr.String(), err)
	}

	return nil
}

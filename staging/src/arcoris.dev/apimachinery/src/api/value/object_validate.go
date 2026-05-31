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

package value

// validateObjectMember checks one object member before payload insertion.
//
// The existing slice contains only already validated and cloned members. Passing
// it in keeps duplicate-name detection local to object construction without
// storing an index in the final payload.
func validateObjectMember(index int, member Member, existing []Member) error {
	if member.Name == "" {
		return errorf(
			objectMemberNamePath(index),
			ErrEmptyName,
			ErrorReasonEmptyName,
			"object member name is empty",
		)
	}

	if member.Value.IsZero() {
		return errorf(
			objectMemberValuePath(index),
			ErrInvalidMember,
			ErrorReasonInvalidValue,
			"object member %q has an invalid zero value",
			member.Name,
		)
	}

	if hasObjectMemberName(existing, member.Name) {
		return errorf(
			objectMemberNamePath(index),
			ErrDuplicateName,
			ErrorReasonDuplicateName,
			"object member name %q is duplicated",
			member.Name,
		)
	}

	return nil
}

// hasObjectMemberName performs the intentionally small linear duplicate check.
//
// It trades O(n) lookup for lower allocation and simpler payload invariants,
// which is the better default for short API objects.
func hasObjectMemberName(members []Member, name string) bool {
	for _, member := range members {
		if member.Name == name {
			return true
		}
	}

	return false
}

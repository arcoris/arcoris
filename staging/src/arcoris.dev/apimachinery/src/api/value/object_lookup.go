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

// findObjectMember returns the ordered member index for name.
//
// Object payloads intentionally avoid storing name indexes. API payload objects
// are expected to be short, and a linear scan keeps construction, cloning, and
// view creation allocation-light. A negative result means the member is absent.
func findObjectMember(members []Member, name string) int {
	for i, member := range members {
		if member.Name == name {
			return i
		}
	}

	return -1
}

// hasObjectMemberName reports whether members already contains name.
//
// Object construction uses this to reject duplicate names before appending the
// new member. Sharing the lookup primitive with ObjectView keeps duplicate
// detection and read access on the same linear semantics.
func hasObjectMemberName(members []Member, name string) bool {
	return findObjectMember(members, name) >= 0
}

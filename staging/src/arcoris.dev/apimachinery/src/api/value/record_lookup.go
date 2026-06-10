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

// findRecordMember returns the ordered member index for name.
//
// Record payloads intentionally avoid storing name indexes. API payload records
// are expected to be short, and a linear scan keeps stored payloads compact. A
// negative result means the member is absent.
func findRecordMember(members []RecordMember, name MemberName) int {
	for i, member := range members {
		if member.Name == name {
			return i
		}
	}

	return -1
}

// hasRecordMemberName reports whether members already contains name.
func hasRecordMemberName(members []RecordMember, name MemberName) bool {
	return findRecordMember(members, name) >= 0
}

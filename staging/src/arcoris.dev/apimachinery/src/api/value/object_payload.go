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

// objectPayload stores object members in caller order.
//
// The payload intentionally has no name index. Objects are expected to be small,
// and linear lookup avoids extra allocations and duplicate invariants during
// construction, cloning, and view creation.
type objectPayload struct {
	// members contains cloned member values in stable caller order.
	members []Member
}

// newObjectPayload validates members and clones values into caller order.
//
// Duplicate detection scans the members already accepted into the payload. This
// keeps the constructor allocation profile small and makes member order the only
// stored object invariant.
func newObjectPayload(members []Member) (objectPayload, error) {
	if len(members) == 0 {
		return objectPayload{}, nil
	}

	payload := objectPayload{
		members: make([]Member, 0, len(members)),
	}

	for i, member := range members {
		if err := validateObjectMember(i, member, payload.members); err != nil {
			return objectPayload{}, err
		}

		payload.members = append(payload.members, ObjectMember(member.Name, member.Value))
	}

	return payload, nil
}

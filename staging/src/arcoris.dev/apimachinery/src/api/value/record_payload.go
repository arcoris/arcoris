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

const recordDuplicateMapThreshold = 16

// recordPayload stores record members in caller order.
//
// The payload intentionally has no name index. Records are expected to be small,
// and linear lookup keeps stored payloads compact. Large construction inputs use
// a temporary map only while checking duplicate names.
type recordPayload struct {
	// members contains cloned member values in stable caller order.
	members []RecordMember
}

// newRecordPayload validates members and clones values into caller order.
//
// Member values are cloned exactly once at this construction boundary. Duplicate
// detection is linear for small records and map-backed for larger inputs; no
// lookup structure is retained in the payload.
func newRecordPayload(members []RecordMember) (recordPayload, error) {
	if len(members) == 0 {
		return recordPayload{}, nil
	}

	payload := recordPayload{
		members: make([]RecordMember, 0, len(members)),
	}

	if len(members) > recordDuplicateMapThreshold {
		seen := make(map[MemberName]int, len(members))
		for i, member := range members {
			if err := validateRecordMember(i, member); err != nil {
				return recordPayload{}, err
			}
			if first, ok := seen[member.Name]; ok {
				return recordPayload{}, recordDuplicateMemberName(i, member.Name, first)
			}

			seen[member.Name] = i
			payload.members = append(payload.members, cloneRecordMember(member))
		}

		return payload, nil
	}

	for i, member := range members {
		if err := validateRecordMember(i, member); err != nil {
			return recordPayload{}, err
		}
		if first := findRecordMember(payload.members, member.Name); first >= 0 {
			return recordPayload{}, recordDuplicateMemberName(i, member.Name, first)
		}

		payload.members = append(payload.members, cloneRecordMember(member))
	}

	return payload, nil
}

func cloneRecordMember(member RecordMember) RecordMember {
	return RecordMember{Name: member.Name, Value: member.Value.Clone()}
}

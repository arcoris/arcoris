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

// mapPayload stores map entries in caller order.
//
// The payload intentionally has no key index. Maps in this value model are
// expected to be small, and linear lookup avoids extra allocation in
// construction, cloning, and view creation.
type mapPayload struct {
	// entries contains cloned entry values in stable caller order.
	entries []Entry
}

// newMapPayload validates entries and clones values into caller order.
//
// Duplicate detection scans entries already accepted into the payload. This
// keeps entry order as the only stored map invariant.
func newMapPayload(entries []Entry) (mapPayload, error) {
	payload := mapPayload{
		entries: make([]Entry, 0, len(entries)),
	}

	for i, entry := range entries {
		if err := validateMapEntry(i, entry, payload.entries); err != nil {
			return mapPayload{}, err
		}

		payload.entries = append(payload.entries, MapEntry(entry.Key, entry.Value))
	}

	return payload.compact(), nil
}

// compact removes empty backing storage from empty maps.
//
// Empty maps are valid, but nil storage avoids retaining an unused backing
// array.
func (p mapPayload) compact() mapPayload {
	if len(p.entries) == 0 {
		return mapPayload{}
	}

	return p
}

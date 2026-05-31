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

// objectPayload stores object fields in caller order.
//
// The payload intentionally has no name index. Objects are expected to be small,
// and linear lookup avoids extra allocations and duplicate invariants during
// construction, cloning, and view creation.
type objectPayload struct {
	// fields contains cloned field values in stable caller order.
	fields []Field
}

// newObjectPayload validates fields and clones values into caller order.
//
// Duplicate detection scans the fields already accepted into the payload. This
// keeps the constructor allocation profile small and makes field order the only
// stored object invariant.
func newObjectPayload(fields []Field) (objectPayload, error) {
	payload := objectPayload{
		fields: make([]Field, 0, len(fields)),
	}

	for i, field := range fields {
		if err := validateObjectField(i, field, payload.fields); err != nil {
			return objectPayload{}, err
		}

		payload.fields = append(payload.fields, ObjectField(field.Name, field.Value))
	}

	return payload.compact(), nil
}

// compact removes empty backing storage from empty objects.
//
// Empty objects are valid, but keeping nil storage for them avoids retaining an
// otherwise unused zero-length backing array.
func (p objectPayload) compact() objectPayload {
	if len(p.fields) == 0 {
		return objectPayload{}
	}

	return p
}

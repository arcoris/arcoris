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

package capacity

// CheckedAdd returns v+other and false on amount overflow.
func (v Vector) CheckedAdd(other Vector) (Vector, bool) {
	result := make([]Entry, 0, len(v.entries)+len(other.entries))

	i, j := 0, 0
	for i < len(v.entries) || j < len(other.entries) {
		leftRemaining := i < len(v.entries)
		rightRemaining := j < len(other.entries)
		takeLeft := leftRemaining && (!rightRemaining || v.entries[i].Resource < other.entries[j].Resource)
		takeRight := rightRemaining && (!leftRemaining || other.entries[j].Resource < v.entries[i].Resource)

		switch {
		case takeLeft:
			result = append(result, v.entries[i])
			i++

		case takeRight:
			result = append(result, other.entries[j])
			j++

		default:
			sum, ok := v.entries[i].Amount.CheckedAdd(other.entries[j].Amount)
			if !ok {
				return Vector{}, false
			}

			result = append(result, Entry{Resource: v.entries[i].Resource, Amount: sum})
			i++
			j++
		}
	}

	return Vector{entries: result}, true
}

// CheckedSub returns v-other and false when v does not cover other.
func (v Vector) CheckedSub(other Vector) (Vector, bool) {
	result := make([]Entry, 0, len(v.entries))

	i, j := 0, 0
	for i < len(v.entries) {
		entry := v.entries[i]
		otherExhausted := j >= len(other.entries)
		entryOnlyInReceiver := !otherExhausted && entry.Resource < other.entries[j].Resource

		if otherExhausted || entryOnlyInReceiver {
			result = append(result, entry)
			i++
			continue
		}

		if other.entries[j].Resource < entry.Resource {
			return Vector{}, false
		}

		diff, ok := entry.Amount.CheckedSub(other.entries[j].Amount)
		if !ok {
			return Vector{}, false
		}

		if diff.IsPositive() {
			result = append(result, Entry{Resource: entry.Resource, Amount: diff})
		}

		i++
		j++
	}

	if j < len(other.entries) {
		return Vector{}, false
	}

	return Vector{entries: result}, true
}

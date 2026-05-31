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

// validateMapEntry checks one map entry before payload insertion.
//
// The existing slice contains only already validated and cloned entries. Passing
// it in keeps duplicate-key detection local to map construction without storing
// an index in the final payload.
func validateMapEntry(index int, entry Entry, existing []Entry) error {
	if entry.Key == "" {
		return errorf(
			mapEntryKeyPath(index),
			ErrEmptyKey,
			ErrorReasonEmptyKey,
			"map entry key is empty",
		)
	}

	if entry.Value.IsZero() {
		return errorf(
			mapEntryValuePath(index),
			ErrInvalidEntry,
			ErrorReasonInvalidValue,
			"map entry %q has an invalid zero value",
			entry.Key,
		)
	}

	if hasMapEntryKey(existing, entry.Key) {
		return errorf(
			mapEntryKeyPath(index),
			ErrDuplicateKey,
			ErrorReasonDuplicateKey,
			"map entry key %q is duplicated",
			entry.Key,
		)
	}

	return nil
}

// hasMapEntryKey performs the intentionally small linear duplicate check.
//
// It trades O(n) lookup for lower allocation and simpler payload invariants,
// which is the better default for short API maps.
func hasMapEntryKey(entries []Entry, key string) bool {
	for _, entry := range entries {
		if entry.Key == key {
			return true
		}
	}

	return false
}

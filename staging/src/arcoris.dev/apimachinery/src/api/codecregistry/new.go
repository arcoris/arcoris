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

package codecregistry

// New validates configured codec registrations and returns an immutable registry.
//
// New accepts normalizable codec.Info metadata, stores normalized detached
// metadata in registry entries, and rejects invalid or non-normalizable
// metadata. Registered codec implementations are kept as-is, but their metadata
// is snapshotted during construction.
//
// EntryID is the only unique registry identity. Media types and formats are
// grouping attributes, so multiple configured codec instances may share the same
// media type and format.
//
// Construction is all-or-nothing. If any registration is invalid, no partial
// Registry is returned.
func New(registrations ...Registration) (Registry, error) {
	if len(registrations) == 0 {
		return Registry{}, nil
	}

	entries := make([]Entry, 0, len(registrations))
	seenIDs := make(map[EntryID]int, len(registrations))

	for i, current := range registrations {
		entry, err := buildEntry(i, current)
		if err != nil {
			return Registry{}, err
		}
		if err := checkEntryConflicts(i, entry, seenIDs); err != nil {
			return Registry{}, err
		}

		entries = append(entries, entry)
	}

	return buildRegistry(entries), nil
}

/*
   Copyright 2026 The ARCORIS Authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package schema

import (
	"fmt"
	"strings"
)

// GroupVersions is an ordered preference list of GroupVersion values.
//
// The order is authoritative: matching helpers walk this list first, and the
// first preferred GroupVersion that has a candidate wins. Candidate order is
// used only to break ties within the same preferred GroupVersion.
type GroupVersions []GroupVersion

// Validate checks that every listed GroupVersion is a complete canonical identity.
//
// Empty lists are valid and mean no preferences are available. Invalid entries
// are reported with their index so callers can diagnose configuration errors.
func (gvs GroupVersions) Validate() error {
	for i, gv := range gvs {
		if err := gv.Validate(); err != nil {
			return fmt.Errorf("schema: invalid group/version preference at index %d: %w", i, err)
		}
	}
	return nil
}

// IsZero reports whether the preference list is empty.
//
// A zero preference list never matches candidates.
func (gvs GroupVersions) IsZero() bool {
	return len(gvs) == 0
}

// Identifier returns a deterministic string representation of the preference list.
//
// The format preserves list order because order is part of the GroupVersions
// contract. It is intended for diagnostics and stable map keys, not parsing.
func (gvs GroupVersions) Identifier() string {
	if len(gvs) == 0 {
		return "[]"
	}
	entries := make([]string, len(gvs))
	for i, gv := range gvs {
		entries[i] = gv.Identifier()
	}
	return "[" + strings.Join(entries, ", ") + "]"
}

// String returns the same deterministic representation as Identifier.
//
// Keeping String and Identifier identical avoids a second display convention for
// ordered version preferences.
func (gvs GroupVersions) String() string {
	return gvs.Identifier()
}

// Contains reports whether the list contains the exact group/version value.
//
// It performs structural equality and does not validate either side.
func (gvs GroupVersions) Contains(gv GroupVersion) bool {
	for _, candidate := range gvs {
		if candidate == gv {
			return true
		}
	}
	return false
}

// KindForGroupVersionKinds returns the preferred kind candidate.
//
// The receiver's order is authoritative. The first GroupVersion in the
// preference list that matches any candidate wins, regardless of the input
// candidate order.
func (gvs GroupVersions) KindForGroupVersionKinds(kinds []GroupVersionKind) (GroupVersionKind, bool) {
	for _, preferred := range gvs {
		for _, candidate := range kinds {
			if candidate.GroupVersion() == preferred {
				return candidate, true
			}
		}
	}
	return GroupVersionKind{}, false
}

// ResourceForGroupVersionResources returns the preferred resource candidate.
//
// The receiver's order is authoritative. The first GroupVersion in the
// preference list that matches any candidate wins, regardless of the input
// candidate order.
func (gvs GroupVersions) ResourceForGroupVersionResources(resources []GroupVersionResource) (GroupVersionResource, bool) {
	for _, preferred := range gvs {
		for _, candidate := range resources {
			if candidate.GroupVersion() == preferred {
				return candidate, true
			}
		}
	}
	return GroupVersionResource{}, false
}

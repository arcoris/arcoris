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

package objectquery

import "arcoris.dev/apimachinery/api/objectstore"

// Predicate is a validated, canonical object list item predicate.
type Predicate struct {
	// identity is the canonical identity predicate.
	identity IdentitySelector

	// labels is the canonical label predicate.
	labels LabelSelector

	// annotations is the canonical annotation predicate.
	annotations AnnotationSelector
}

// IsZero reports whether p matches every item.
func (p Predicate) IsZero() bool {
	return p.identity.IsZero() && p.labels.IsZero() && p.annotations.IsZero()
}

// Match reports whether item satisfies every predicate section.
func (p Predicate) Match(item objectstore.ListItem) bool {
	return p.identity.match(item) &&
		p.labels.match(item) &&
		p.annotations.match(item)
}

// Filter returns items that match p while preserving input order.
//
// Filter does not mutate items and does not clone item state. A nil input slice
// returns nil. An empty non-nil input slice is returned unchanged. Non-empty
// inputs produce a new result slice containing shallow item copies.
func (p Predicate) Filter(items []objectstore.ListItem) []objectstore.ListItem {
	if items == nil {
		return nil
	}
	if len(items) == 0 {
		return items
	}

	out := make([]objectstore.ListItem, 0, len(items))
	for _, item := range items {
		if p.Match(item) {
			out = append(out, item)
		}
	}

	return out
}

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

// Filter returns matching items in input order.
//
// Filter is shallow: it does not clone objectstore.State. Store, lifecycle,
// and cache layers own detachment boundaries.
//
// Nil input stays nil. Empty non-nil input stays non-nil. A non-empty input
// that has no matches returns nil to avoid allocating an empty result.
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
	if len(out) == 0 {
		return nil
	}

	return out
}

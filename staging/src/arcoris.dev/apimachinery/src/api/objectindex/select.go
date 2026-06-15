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

package objectindex

import (
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

// Select compiles query and returns matching items in original input order.
//
// Invalid objectquery values are returned unchanged from objectquery. Empty
// indexes and no-match selections return nil. Returned items are shallow copies
// and object states are not cloned.
func (idx Index) Select(query objectquery.Query) ([]objectstore.ListItem, error) {
	predicate, err := objectquery.Compile(query)
	if err != nil {
		return nil, err
	}

	return idx.SelectPredicate(predicate), nil
}

// SelectPredicate returns matching items for a compiled objectquery predicate.
//
// The predicate's detached canonical query is used only for candidate
// narrowing. The predicate itself is always applied before a ListItem is
// returned, preserving objectquery full-scan semantics.
func (idx Index) SelectPredicate(predicate objectquery.Predicate) []objectstore.ListItem {
	if len(idx.items) == 0 {
		return nil
	}
	if predicate.IsZero() {
		return append([]objectstore.ListItem(nil), idx.items...)
	}

	return idx.selectFromCandidates(predicate, idx.planCandidates(predicate))
}

// selectFromCandidates scans the original item order and applies predicate to
// every candidate position.
//
// Returning only after predicate.Match keeps objectindex equivalent to
// objectquery full-scan filtering even when planning used only a subset of the
// query.
func (idx Index) selectFromCandidates(
	predicate objectquery.Predicate,
	plan candidatePlan,
) []objectstore.ListItem {
	out := make([]objectstore.ListItem, 0, len(idx.items))
	for pos, item := range idx.items {
		if !plan.includes(pos) {
			continue
		}
		if predicate.Match(item) {
			out = append(out, item)
		}
	}
	if len(out) == 0 {
		return nil
	}

	return out
}

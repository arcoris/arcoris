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

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

// TestPredicateFilterNilEmptyNoMatchAndOrder verifies Filter's shallow ordered
// result contract.
func TestPredicateFilterNilEmptyNoMatchAndOrder(t *testing.T) {
	predicate := mustPredicate(t, mustQ(LabelEquals("env", "prod")))

	if got := predicate.Filter(nil); got != nil {
		t.Fatalf("Filter(nil) = %#v; want nil", got)
	}

	empty := []objectstore.ListItem{}
	if got := predicate.Filter(empty); got == nil {
		t.Fatal("Filter(empty) = nil; want empty non-nil")
	}

	requireNames(t, predicate.Filter(testItems()), "worker-1", "worker-3")
	requireNames(t, mustPredicate(t, None()).Filter(testItems()))
}

// BenchmarkPredicateFilter100 covers the common full-scan path without making
// planning or cache internals part of objectquery.
func BenchmarkPredicateFilter100(b *testing.B) {
	items := make([]objectstore.ListItem, 0, 100)
	for i := 0; i < 25; i++ {
		items = append(items, testItems()...)
	}
	predicate := mustPredicate(b, mustQ(LabelEquals("env", "prod")))

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = predicate.Filter(items)
	}
}

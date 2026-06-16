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

package objectcache

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
)

func TestCollectionListNoMatchesReturnsNil(t *testing.T) {
	col := mustCollection(t, testListResult(8, testItems()...))
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "missing")),
	}
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	if got := col.list(predicate); got != nil {
		t.Fatalf("list() = %#v; want nil", got)
	}
}

func TestCollectionListUsesPredicateAfterIndexPlan(t *testing.T) {
	col := mustCollection(t, testListResult(8, testItems()...))
	query := objectquery.Query{
		Labels: mustLabelSelector(
			t,
			mustLabelExists(t, "env"),
			mustLabelNotEquals(t, "tier", "frontend"),
		),
	}

	assertCollectionListEquivalent(t, col, query)
}

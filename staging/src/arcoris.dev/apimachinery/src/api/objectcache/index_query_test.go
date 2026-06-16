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

func TestIndexesPlanIdentityLabelAndAnnotation(t *testing.T) {
	items := testItems()
	idx := testIndexes()
	query := objectquery.Query{
		Identity: objectquery.IdentitySelector{
			Namespace: mustNamespaceEquals(t, "system"),
		},
		Labels:      mustLabelSelector(t, mustLabelIn(t, "env", "prod", "qa")),
		Annotations: mustAnnotationSelector(t, mustAnnotationEquals(t, "team", "core")),
	}
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	plan := idx.plan(predicate)

	if !plan.includes(items[0].Key) {
		t.Fatal("plan excludes matching worker-1 candidate")
	}
	if plan.includes(items[1].Key) {
		t.Fatal("plan includes worker-2; annotation team should narrow it out")
	}
	if plan.includes(items[2].Key) {
		t.Fatal("plan includes worker-3; namespace should narrow it out")
	}
}

func TestIndexesPlanNegativeOnlyIsUnconstrained(t *testing.T) {
	idx := newIndexes()
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelDoesNotExist(t, "missing")),
	}
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	plan := idx.plan(predicate)

	if plan.constrained {
		t.Fatal("plan.constrained = true; want false for residual-only query")
	}
}

func testIndexes() indexes {
	idx := newIndexes()
	for _, item := range testItems() {
		idx.add(item)
	}

	return idx
}

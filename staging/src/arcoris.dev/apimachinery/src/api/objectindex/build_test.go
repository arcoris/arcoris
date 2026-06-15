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
	"testing"

	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

func TestBuildNil(t *testing.T) {
	idx := Build(nil)
	if !idx.IsZero() {
		t.Fatal("IsZero() = false; want true")
	}
	if idx.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", idx.Len())
	}

	got, err := idx.Select(objectquery.Query{})
	requireNoError(t, err)
	if got != nil {
		t.Fatalf("Select() = %#v; want nil", got)
	}
}

func TestBuildEmpty(t *testing.T) {
	idx := Build([]objectstore.ListItem{})
	if !idx.IsZero() {
		t.Fatal("IsZero() = false; want true")
	}
	if idx.Len() != 0 {
		t.Fatalf("Len() = %d; want 0", idx.Len())
	}
	if got := idx.SelectPredicate(objectquery.Predicate{}); got != nil {
		t.Fatalf("SelectPredicate() = %#v; want nil", got)
	}
}

func TestBuildCopiesItemSlice(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 2, labelsMap("env", "qa"), nil),
	}
	idx := Build(items)
	items[0] = testItem("system", "changed", 99, nil, nil)

	got := idx.SelectPredicate(objectquery.Predicate{})
	requireItemOrder(t, got, itemRef{"system", "worker-1", 1}, itemRef{"system", "worker-2", 2})

	got[0] = testItem("system", "mutated-result", 100, nil, nil)
	again := idx.SelectPredicate(objectquery.Predicate{})
	requireItemOrder(t, again, itemRef{"system", "worker-1", 1}, itemRef{"system", "worker-2", 2})
}

func TestBuildDoesNotCloneState(t *testing.T) {
	idx := Build([]objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
	})

	got := idx.SelectPredicate(objectquery.Predicate{})
	got[0].State.Object.ObjectMeta.Labels[labels.Key("env")] = labels.Value("mutated")

	again := idx.SelectPredicate(objectquery.Predicate{})
	value, ok := again[0].State.Object.ObjectMeta.Labels.Get("env")
	if !ok || value != "mutated" {
		t.Fatalf("label env = %q, %v; want mutated, true", value, ok)
	}
}

func TestBuildPreservesDuplicateObjectKeys(t *testing.T) {
	items := []objectstore.ListItem{
		testItem("system", "worker-1", 1, labelsMap("env", "prod"), nil),
		testItem("system", "worker-1", 2, labelsMap("env", "prod"), nil),
		testItem("system", "worker-2", 3, labelsMap("env", "qa"), nil),
	}
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelEquals(t, "env", "prod")),
	}

	assertIndexMatchesFullScan(t, items, query)

	got, err := Build(items).Select(query)
	requireNoError(t, err)
	requireItemOrder(t, got, itemRef{"system", "worker-1", 1}, itemRef{"system", "worker-1", 2})
}

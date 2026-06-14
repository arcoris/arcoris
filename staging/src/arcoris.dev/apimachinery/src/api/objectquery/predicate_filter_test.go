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
	"arcoris.dev/apimachinery/api/value"
)

func TestPredicateFilterNilInputReturnsNil(t *testing.T) {
	predicate, err := Compile(Query{})
	requireNoError(t, err)

	if got := predicate.Filter(nil); got != nil {
		t.Fatalf("Filter(nil) = %#v; want nil", got)
	}
}

func TestPredicateFilterEmptyInputReturnsEmptyNonNil(t *testing.T) {
	predicate, err := Compile(Query{})
	requireNoError(t, err)
	input := []objectListItem{}
	items := objectListItems(input)

	got := predicate.Filter(items)
	if got == nil {
		t.Fatal("Filter(empty) returned nil")
	}
	if len(got) != 0 {
		t.Fatalf("len = %d; want 0", len(got))
	}
}

func TestPredicateFilterPreservesOrderAndSelectsMatches(t *testing.T) {
	selector := mustLabelSelector(t, mustLabelEquals(t, "env", "prod"))
	predicate, err := Compile(Query{Labels: selector})
	requireNoError(t, err)
	items := objectListItems([]objectListItem{
		{namespace: "system", name: "first", labels: map[string]string{"env": "prod"}},
		{namespace: "system", name: "second", labels: map[string]string{"env": "qa"}},
		{namespace: "system", name: "third", labels: map[string]string{"env": "prod"}},
	})

	got := predicate.Filter(items)
	if len(got) != 2 {
		t.Fatalf("len = %d; want 2", len(got))
	}
	if got[0].Key.Object.Name != "first" || got[1].Key.Object.Name != "third" {
		t.Fatalf("order = %s, %s; want first, third", got[0].Key.Object.Name, got[1].Key.Object.Name)
	}
}

func TestPredicateFilterZeroPredicateReturnsAllItems(t *testing.T) {
	predicate, err := Compile(Query{})
	requireNoError(t, err)
	items := objectListItems([]objectListItem{
		{namespace: "system", name: "first"},
		{namespace: "system", name: "second"},
	})

	got := predicate.Filter(items)
	if len(got) != len(items) {
		t.Fatalf("len = %d; want %d", len(got), len(items))
	}
}

func TestPredicateFilterDoesNotMutateInputAndDoesNotCloneState(t *testing.T) {
	predicate, err := Compile(Query{})
	requireNoError(t, err)
	observed := value.StringValue("observed")
	items := objectListItems([]objectListItem{{namespace: "system", name: "worker"}})
	items[0].State.Object = items[0].State.Object.WithObserved(observed)

	got := predicate.Filter(items)
	if len(got) != 1 {
		t.Fatalf("len = %d; want 1", len(got))
	}
	if got[0].State.Object.Observed != items[0].State.Object.Observed {
		t.Fatal("Filter cloned state; observed pointer differs")
	}

	got[0].State.Object.Desired = value.StringValue("mutated")
	if text, _ := items[0].State.Object.Desired.AsString(); text != "worker" {
		t.Fatalf("input desired = %q; want worker", text)
	}
}

type objectListItem struct {
	namespace   string
	name        string
	labels      map[string]string
	annotations map[string]string
}

func objectListItems(items []objectListItem) []objectstore.ListItem {
	out := make([]objectstore.ListItem, 0, len(items))
	for _, item := range items {
		out = append(out, testItem(item.namespace, item.name, item.labels, item.annotations))
	}

	return out
}

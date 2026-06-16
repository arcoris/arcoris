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
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false; want true", err, target)
	}
}

func requireErrorNotIs(t *testing.T, err error, target error) {
	t.Helper()

	if errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = true; want false", err, target)
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireItemOrder(t *testing.T, got []objectstore.ListItem, want ...itemRef) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len = %d; want %d; got=%v want=%v", len(got), len(want), itemRefs(got), want)
	}
	for i, ref := range want {
		if got[i].Key.Object.Namespace != ref.namespace ||
			got[i].Key.Object.Name != ref.name ||
			got[i].State.Revision != ref.revision {
			t.Fatalf("item[%d] = %v; want %v", i, itemRefs(got), want)
		}
	}
}

func requireSameItems(t *testing.T, got []objectstore.ListItem, want []objectstore.ListItem) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("len = %d; want %d; got=%v want=%v", len(got), len(want), itemRefs(got), itemRefs(want))
	}
	for i := range want {
		if !got[i].Key.Equal(want[i].Key) || got[i].State.Revision != want[i].State.Revision {
			t.Fatalf("item[%d] = %v; want %v", i, itemRefs(got), itemRefs(want))
		}
	}
}

func itemRefs(items []objectstore.ListItem) []itemRef {
	if items == nil {
		return nil
	}

	out := make([]itemRef, len(items))
	for i, item := range items {
		out[i] = itemRef{
			namespace: item.Key.Object.Namespace,
			name:      item.Key.Object.Name,
			revision:  item.State.Revision,
		}
	}

	return out
}

func desiredString(t *testing.T, item objectstore.ListItem) string {
	t.Helper()

	got, ok := item.State.Object.Desired.AsString()
	if !ok {
		t.Fatalf("Desired is not a string: %#v", item.State.Object.Desired)
	}

	return got
}

func labelValue(t *testing.T, item objectstore.ListItem, key string) string {
	t.Helper()

	got, ok := item.State.Object.ObjectMeta.Labels.Get(labels.Key(key))
	if !ok {
		t.Fatalf("label %q missing", key)
	}

	return got.String()
}

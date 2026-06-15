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

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

const testGroup = apiidentity.Group("control.arcoris.dev")

type itemRef struct {
	namespace metaidentity.Namespace
	name      metaidentity.Name
	revision  objectstore.Revision
}

func TestIndexLenAndIsZero(t *testing.T) {
	tests := []struct {
		name   string
		items  []objectstore.ListItem
		len    int
		isZero bool
	}{
		{
			name:   "nil",
			items:  nil,
			len:    0,
			isZero: true,
		},
		{
			name:   "empty",
			items:  []objectstore.ListItem{},
			len:    0,
			isZero: true,
		},
		{
			name:   "items",
			items:  testItems(),
			len:    len(testItems()),
			isZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := Build(tt.items)
			if got := idx.Len(); got != tt.len {
				t.Fatalf("Len() = %d; want %d", got, tt.len)
			}
			if got := idx.IsZero(); got != tt.isZero {
				t.Fatalf("IsZero() = %v; want %v", got, tt.isZero)
			}
		})
	}
}

func testItems() []objectstore.ListItem {
	return []objectstore.ListItem{
		testItem(
			"system",
			"worker-1",
			1,
			map[string]string{"env": "prod", "tier": "backend"},
			map[string]string{"team": "core", "zone": "east"},
		),
		testItem(
			"system",
			"worker-2",
			2,
			map[string]string{"env": "qa", "tier": "backend"},
			map[string]string{"team": "tools", "zone": "west"},
		),
		testItem(
			"other",
			"worker-3",
			3,
			map[string]string{"env": "prod", "tier": "frontend"},
			map[string]string{"team": "core", "zone": "west"},
		),
		testItem(
			"",
			"cluster-worker",
			4,
			map[string]string{"env": "prod"},
			map[string]string{"team": "platform"},
		),
		testItem("system", "worker-4", 5, nil, nil),
	}
}

func testGVR() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    testGroup,
		Version:  "v1",
		Resource: "workers",
	}
}

func testGVK() apiidentity.GroupVersionKind {
	return apiidentity.GroupVersionKind{
		Group:   testGroup,
		Version: "v1",
		Kind:    "Worker",
	}
}

func testItem(
	namespace metaidentity.Namespace,
	name metaidentity.Name,
	revision objectstore.Revision,
	rawLabels map[string]string,
	rawAnnotations map[string]string,
) objectstore.ListItem {
	objectName := metaidentity.ObjectName{Namespace: namespace, Name: name}
	objectMeta := meta.ObjectMeta{
		Name:        name,
		Namespace:   namespace,
		Labels:      mustLabels(rawLabels),
		Annotations: mustAnnotations(rawAnnotations),
	}
	obj := object.New[value.Value, value.Value](
		meta.FromGroupVersionKind(testGVK()),
		objectMeta,
		value.StringValue(name.String()),
	)

	return objectstore.ListItem{
		Key: objectstore.MustKey(testGVR(), objectName),
		State: objectstore.State{
			Object:    obj,
			Ownership: objectownership.EmptyState(),
			Revision:  revision,
		},
	}
}

func labelsMap(key string, val string) map[string]string {
	return map[string]string{key: val}
}

func annotationsMap(key string, val string) map[string]string {
	return map[string]string{key: val}
}

func mustLabels(values map[string]string) labels.Set {
	set, err := labels.FromStrings(values)
	if err != nil {
		panic(err)
	}

	return set
}

func mustAnnotations(values map[string]string) annotations.Set {
	set, err := annotations.FromStrings(values)
	if err != nil {
		panic(err)
	}

	return set
}

func assertIndexMatchesFullScan(
	t *testing.T,
	items []objectstore.ListItem,
	query objectquery.Query,
) {
	t.Helper()

	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	want := predicate.Filter(items)
	idx := Build(items)
	got, err := idx.Select(query)
	requireNoError(t, err)
	requireSameItems(t, got, want)

	gotPredicate := idx.SelectPredicate(predicate)
	requireSameItems(t, gotPredicate, want)
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

func requirePositions(t *testing.T, got []int, want ...int) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("positions = %v; want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("positions = %v; want %v", got, want)
		}
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func mustNamespaceEquals(t *testing.T, namespace metaidentity.Namespace) objectquery.NamespaceRequirement {
	t.Helper()

	req, err := objectquery.NamespaceEquals(namespace)
	requireNoError(t, err)

	return req
}

func mustNameEquals(t *testing.T, name metaidentity.Name) objectquery.NameRequirement {
	t.Helper()

	req, err := objectquery.NameEquals(name)
	requireNoError(t, err)

	return req
}

func mustLabelExists(t *testing.T, key string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelExists(key)
	requireNoError(t, err)

	return req
}

func mustLabelDoesNotExist(t *testing.T, key string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustLabelEquals(t *testing.T, key string, val string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustLabelNotEquals(t *testing.T, key string, val string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelNotEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustLabelIn(t *testing.T, key string, values ...string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelNotIn(t *testing.T, key string, values ...string) objectquery.LabelRequirement {
	t.Helper()

	req, err := objectquery.LabelNotIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelSelector(
	t *testing.T,
	requirements ...objectquery.LabelRequirement,
) objectquery.LabelSelector {
	t.Helper()

	selector, err := objectquery.NewLabelSelector(requirements...)
	requireNoError(t, err)

	return selector
}

func mustAnnotationExists(t *testing.T, key string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationExists(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationDoesNotExist(t *testing.T, key string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationEquals(t *testing.T, key string, val string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotEquals(t *testing.T, key string, val string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationNotEquals(key, val)
	requireNoError(t, err)

	return req
}

func mustAnnotationIn(t *testing.T, key string, values ...string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotIn(t *testing.T, key string, values ...string) objectquery.AnnotationRequirement {
	t.Helper()

	req, err := objectquery.AnnotationNotIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationSelector(
	t *testing.T,
	requirements ...objectquery.AnnotationRequirement,
) objectquery.AnnotationSelector {
	t.Helper()

	selector, err := objectquery.NewAnnotationSelector(requirements...)
	requireNoError(t, err)

	return selector
}

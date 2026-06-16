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

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

const testGroup = apiidentity.Group("control.arcoris.dev")

type itemRef struct {
	namespace metaidentity.Namespace
	name      metaidentity.Name
	revision  objectstore.Revision
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
		testItem("system", "worker-4", 4, nil, nil),
	}
}

func testListResult(revision objectstore.Revision, items ...objectstore.ListItem) objectstore.ListResult {
	return objectstore.ListResult{Items: items, Revision: revision}
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

func mustSnapshot(t *testing.T, result objectstore.ListResult) Snapshot {
	t.Helper()

	snapshot, err := NewSnapshot(result)
	requireNoError(t, err)

	return snapshot
}

func mutateItem(item *objectstore.ListItem, desired string) {
	item.State.Object.Desired = value.StringValue(desired)
	if item.State.Object.ObjectMeta.Labels == nil {
		item.State.Object.ObjectMeta.Labels = labels.Set{}
	}
	item.State.Object.ObjectMeta.Labels[labels.Key("env")] = labels.Value(desired)
}

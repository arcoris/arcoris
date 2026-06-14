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
	"errors"
	"reflect"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
)

func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireStrings(t *testing.T, got []string, want ...string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("strings = %#v; want %#v", got, want)
	}
}

func requireRequirement(t *testing.T, got metadataRequirement, key string, op Operator, values ...string) {
	t.Helper()
	if got.key != key {
		t.Fatalf("key = %q; want %q", got.key, key)
	}
	if got.op != op {
		t.Fatalf("operator = %s; want %s", got.op, op)
	}
	requireStrings(t, got.values, values...)
}

func testItem(namespace string, name string, labelValues map[string]string, annotationValues map[string]string) objectstore.ListItem {
	labelSet, err := labels.FromStrings(labelValues)
	if err != nil {
		panic(err)
	}
	annotationSet, err := annotations.FromStrings(annotationValues)
	if err != nil {
		panic(err)
	}

	key := objectstore.MustKey(
		apiidentity.GroupVersionResource{
			Group:    "control.arcoris.dev",
			Version:  "v1",
			Resource: "workers",
		},
		metaidentity.ObjectName{
			Namespace: metaidentity.Namespace(namespace),
			Name:      metaidentity.Name(name),
		},
	)

	return objectstore.ListItem{
		Key: key,
		State: objectstore.State{
			Object: object.New[value.Value, value.Value](
				meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
					Group:   "control.arcoris.dev",
					Version: "v1",
					Kind:    "Worker",
				}),
				meta.ObjectMeta{
					Name:        metaidentity.Name(name),
					Namespace:   metaidentity.Namespace(namespace),
					Labels:      labelSet,
					Annotations: annotationSet,
				},
				value.StringValue(name),
			),
			Revision: 1,
		},
	}
}

func mustLabelExists(t *testing.T, key string) LabelRequirement {
	t.Helper()
	req, err := LabelExists(key)
	requireNoError(t, err)

	return req
}

func mustLabelDoesNotExist(t *testing.T, key string) LabelRequirement {
	t.Helper()
	req, err := LabelDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustLabelEquals(t *testing.T, key string, value string) LabelRequirement {
	t.Helper()
	req, err := LabelEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustLabelNotEquals(t *testing.T, key string, value string) LabelRequirement {
	t.Helper()
	req, err := LabelNotEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustLabelIn(t *testing.T, key string, values ...string) LabelRequirement {
	t.Helper()
	req, err := LabelIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelNotIn(t *testing.T, key string, values ...string) LabelRequirement {
	t.Helper()
	req, err := LabelNotIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationExists(t *testing.T, key string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationExists(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationDoesNotExist(t *testing.T, key string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationDoesNotExist(key)
	requireNoError(t, err)

	return req
}

func mustAnnotationEquals(t *testing.T, key string, value string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotEquals(t *testing.T, key string, value string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationNotEquals(key, value)
	requireNoError(t, err)

	return req
}

func mustAnnotationIn(t *testing.T, key string, values ...string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustAnnotationNotIn(t *testing.T, key string, values ...string) AnnotationRequirement {
	t.Helper()
	req, err := AnnotationNotIn(key, values...)
	requireNoError(t, err)

	return req
}

func mustLabelSelector(t *testing.T, requirements ...LabelRequirement) LabelSelector {
	t.Helper()
	selector, err := NewLabelSelector(requirements...)
	requireNoError(t, err)

	return selector
}

func mustAnnotationSelector(t *testing.T, requirements ...AnnotationRequirement) AnnotationSelector {
	t.Helper()
	selector, err := NewAnnotationSelector(requirements...)
	requireNoError(t, err)

	return selector
}

func mustInNamespace(t *testing.T, namespace metaidentity.Namespace) IdentitySelector {
	t.Helper()
	selector, err := InNamespace(namespace)
	requireNoError(t, err)

	return selector
}

func mustWithName(t *testing.T, name metaidentity.Name) IdentitySelector {
	t.Helper()
	selector, err := WithName(name)
	requireNoError(t, err)

	return selector
}

func mustWithObject(t *testing.T, namespace metaidentity.Namespace, name metaidentity.Name) IdentitySelector {
	t.Helper()
	selector, err := WithObject(namespace, name)
	requireNoError(t, err)

	return selector
}

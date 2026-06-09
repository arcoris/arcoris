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

package objectlifecycle

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectmemorystore"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/resourcecatalog"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

const testGroup = apiidentity.Group("control.arcoris.dev")

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

func testName(index int) metaidentity.ObjectName {
	return metaidentity.ObjectName{
		Namespace: "system",
		Name:      metaidentity.Name(fmt.Sprintf("worker-%d", index)),
	}
}

func testTypeMeta() meta.TypeMeta {
	return meta.FromGroupVersionKind(testGVK())
}

func testObjectMeta(index int) meta.ObjectMeta {
	name := testName(index)

	return meta.ObjectMeta{
		Name:      name.Name,
		Namespace: name.Namespace,
	}
}

func desiredDescriptor() types.Descriptor {
	return types.Object(
		types.Field("image").String().Optional(),
		types.Field("replicas").String().Optional(),
	).Descriptor()
}

func observedDescriptor() types.Descriptor {
	return types.Object(
		types.Field("ready").String().Optional(),
	).Descriptor()
}

func testDefinition(opts ...resource.VersionOption) resource.Definition {
	versionOpts := append([]resource.VersionOption{resource.Exposed(), resource.Canonical()}, opts...)

	return resource.NewDefinition(
		testGroup,
		"Worker",
		"workers",
		resource.ScopeNamespaced,
		resource.NewVersion("v1", desiredDescriptor(), versionOpts...),
	)
}

func testCatalog(t *testing.T, definitions ...resource.Definition) *resourcecatalog.Catalog {
	t.Helper()

	catalog := resourcecatalog.New(nil)
	if len(definitions) == 0 {
		definitions = []resource.Definition{testDefinition()}
	}
	for _, definition := range definitions {
		requireNoError(t, catalog.Register(definition))
	}

	return catalog
}

func testStore(t *testing.T) objectstore.Store {
	t.Helper()

	store, err := objectmemorystore.New()
	requireNoError(t, err)

	return store
}

func testExecutor(t *testing.T, opts ...Option) *Executor {
	t.Helper()

	options := []Option{
		WithStore(testStore(t)),
		WithResourceResolver(testCatalog(t)),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
		WithObservedValidator(valuevalidation.SurfaceValidator{}),
	}
	options = append(options, opts...)

	executor, err := NewExecutor(options...)
	requireNoError(t, err)

	return executor
}

func testObject(index int, image string) object.Object[value.Value, value.Value] {
	return object.New[value.Value, value.Value](
		testTypeMeta(),
		testObjectMeta(index),
		objectValue(member("image", value.StringValue(image))),
	)
}

func testObjectWithDesired(index int, desired value.Value) object.Object[value.Value, value.Value] {
	return object.New[value.Value, value.Value](
		testTypeMeta(),
		testObjectMeta(index),
		desired,
	)
}

func testObservedObject(index int, image string, ready string) object.Object[value.Value, value.Value] {
	return object.NewObserved[value.Value, value.Value](
		testTypeMeta(),
		testObjectMeta(index),
		objectValue(member("image", value.StringValue(image))),
		objectValue(member("ready", value.StringValue(ready))),
	)
}

func objectValue(members ...value.Member) value.Value {
	return value.MustObjectValue(members...)
}

func member(name string, val value.Value) value.Member {
	return value.ObjectMember(name, val)
}

func owner(name string) fieldownership.Owner {
	return fieldownership.Owner(name)
}

func createObject(t *testing.T, executor *Executor, index int, image string, owner fieldownership.Owner) Result {
	t.Helper()

	result, err := executor.Create(
		context.Background(),
		CreateRequest{Object: testObject(index, image), Owner: owner},
	)
	requireNoError(t, err)

	return result
}

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

func requireLifecycleError(t *testing.T, err error, target error, reason ErrorReason) *Error {
	t.Helper()

	requireErrorIs(t, err, target)

	var lifecycleErr *Error
	if !errors.As(err, &lifecycleErr) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if lifecycleErr.Reason != reason {
		t.Fatalf("Reason = %q; want %q", lifecycleErr.Reason, reason)
	}

	return lifecycleErr
}

func requireEffect(t *testing.T, result Result, op Operation, effect Effect) {
	t.Helper()
	if result.Operation != op || result.Effect != effect {
		t.Fatalf("result = %#v; want operation=%s effect=%s", result, op, effect)
	}
}

func requireImage(t *testing.T, state objectstore.State, want string) {
	t.Helper()

	objectView, ok := state.Object.Desired.Object()
	if !ok {
		t.Fatalf("desired is not object")
	}
	got, ok := objectView.Get("image")
	if !ok {
		t.Fatalf("desired.image missing")
	}
	text, ok := got.String()
	if !ok || text != want {
		t.Fatalf("desired.image = %q, %v; want %q, true", text, ok, want)
	}
}

func requireOwnedPath(t *testing.T, doc objectownership.Document, owner fieldownership.Owner, path objectownership.Path) {
	t.Helper()

	for _, entry := range doc.Desired.Entries {
		if entry.Owner != owner {
			continue
		}
		for _, field := range entry.Fields {
			if field == path {
				return
			}
		}
	}

	t.Fatalf("ownership path %s for owner %s not found in %#v", path, owner, doc)
}

func ownershipPath(path fieldpath.Path) objectownership.Path {
	return objectownership.Path(path.String())
}

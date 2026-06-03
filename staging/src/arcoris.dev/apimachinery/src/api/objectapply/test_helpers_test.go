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

package objectapply

import (
	"errors"
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/stamp"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/resource"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// testGVK returns the canonical Worker GVK used by objectapply tests.
func testGVK(version apiidentity.Version) apiidentity.GroupVersionKind {
	return apiidentity.GroupVersionKind{
		Group:   "control.arcoris.dev",
		Version: version,
		Kind:    "Worker",
	}
}

// testTypeMeta converts the canonical Worker GVK into object TypeMeta.
func testTypeMeta(version apiidentity.Version) meta.TypeMeta {
	return meta.FromGroupVersionKind(testGVK(version))
}

// testObjectMeta returns live metadata with server-owned fields populated.
func testObjectMeta() meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:            metaidentity.Name("worker"),
		Namespace:       metaidentity.Namespace("system"),
		UID:             metaidentity.UID("uid-1"),
		ResourceVersion: stamp.ResourceVersion("rv-live"),
		Generation:      stamp.Generation(7),
	}
}

// appliedObjectMeta returns applied metadata that exactly identifies live.
func appliedObjectMeta() meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:      metaidentity.Name("worker"),
		Namespace: metaidentity.Namespace("system"),
		UID:       metaidentity.UID("uid-1"),
	}
}

// minimalAppliedObjectMeta returns applied identity without UID.
func minimalAppliedObjectMeta() meta.ObjectMeta {
	return meta.ObjectMeta{
		Name:      metaidentity.Name("worker"),
		Namespace: metaidentity.Namespace("system"),
	}
}

// str builds a string value fixture.
func str(text string) value.Value {
	return value.StringValue(text)
}

// obj builds an object value fixture.
func obj(members ...value.Member) value.Value {
	return value.MustObjectValue(members...)
}

// list builds a list value fixture.
func list(items ...value.Value) value.Value {
	return value.MustListValue(items...)
}

// member builds one object member fixture.
func member(name string, val value.Value) value.Member {
	return value.ObjectMember(name, val)
}

// path parses a test path and panics on invalid literals.
func path(text string) fieldpath.Path {
	p, err := fieldpath.Parse(text)
	if err != nil {
		panic(err)
	}

	return p
}

// fields builds a fieldpath set fixture.
func fields(paths ...fieldpath.Path) fieldpath.Set {
	return fieldpath.MustSet(paths...)
}

// owner builds a field owner fixture.
func owner(name string) fieldownership.Owner {
	return fieldownership.Owner(name)
}

// entry builds one ownership entry fixture.
func entry(name string, paths ...fieldpath.Path) fieldownership.Entry {
	return fieldownership.MustEntry(owner(name), fields(paths...))
}

// desiredOwnership builds object ownership state with Desired ownership only.
func desiredOwnership(entries ...fieldownership.Entry) objectownership.State {
	return objectownership.NewState(fieldownership.MustState(entries...))
}

// desiredDescriptor returns the standard object Desired descriptor.
func desiredDescriptor() types.Type {
	return types.Object(
		types.Field("image").String().Optional(),
		types.Field("replicas").String().Optional(),
	).Type()
}

// observedDescriptor returns the standard observed descriptor.
func observedDescriptor() types.Type {
	return types.Object(
		types.Field("ready").String().Optional(),
	).Type()
}

// mapDesiredDescriptor returns the map Desired descriptor.
func mapDesiredDescriptor() types.Type {
	return types.MapOf(types.String()).Type()
}

// conditionsDescriptor returns the list-map conditions descriptor.
func conditionsDescriptor() types.Type {
	return types.ListOf(
		types.Object(
			types.Field("type").String().Required(),
			types.Field("status").String().Optional(),
		),
	).Map("type").Type()
}

// testResource builds the canonical Worker resource definition.
func testResource(desired types.Type, opts ...resource.VersionOption) resource.Definition {
	options := append([]resource.VersionOption{resource.Exposed(), resource.Canonical()}, opts...)

	return resource.NewDefinition(
		apiidentity.Group("control.arcoris.dev"),
		apiidentity.Kind("Worker"),
		apiidentity.Resource("workers"),
		resource.ScopeNamespaced,
		resource.NewVersion(apiidentity.Version("v1"), desired, options...),
	)
}

// testResourceWithObserved builds a Worker resource that defines Observed.
func testResourceWithObserved(desired types.Type) resource.Definition {
	return testResource(desired, resource.Observed(observedDescriptor()))
}

// testObject builds a live object fixture without Observed.
func testObject(desired value.Value) ValueObject {
	return object.New[value.Value, value.Value](
		testTypeMeta("v1"),
		testObjectMeta(),
		desired,
	)
}

// testObjectObserved builds a live object fixture with Observed.
func testObjectObserved(desired value.Value, observed value.Value) ValueObject {
	return object.NewObserved[value.Value, value.Value](
		testTypeMeta("v1"),
		testObjectMeta(),
		desired,
		observed,
	)
}

// appliedObject builds an applied object fixture without Observed.
func appliedObject(desired value.Value) ValueObject {
	return object.New[value.Value, value.Value](
		testTypeMeta("v1"),
		appliedObjectMeta(),
		desired,
	)
}

// testRequest returns the standard successful apply request.
func testRequest() Request {
	return Request{
		Owner:    owner("user"),
		Live:     testObject(obj(member("image", str("api:v1")), member("replicas", str("3")))),
		Applied:  appliedObject(obj(member("image", str("api:v2")))),
		Resource: testResource(desiredDescriptor()),
	}
}

// readySelector returns the ListMap selector for a Ready condition.
func readySelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry("type", fieldpath.StringLiteral("Ready")),
	)
}

// readyStatusPath returns the semantic path for Ready.status.
func readyStatusPath() fieldpath.Path {
	return fieldpath.RootPath().Select(readySelector()).Field("status")
}

// readyCondition builds a condition item fixture.
func readyCondition(status string) value.Value {
	return obj(member("type", str("Ready")), member("status", str(status)))
}

// requireNoError fails the test when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs fails unless errors.Is matches target.
func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

// requireObjectApplyError checks structured objectapply diagnostics.
func requireObjectApplyError(t *testing.T, err error, path string, reason ErrorReason) {
	t.Helper()

	var applyErr *Error
	if !errors.As(err, &applyErr) {
		t.Fatalf("error type = %T; want *Error", err)
	}
	if applyErr.Path != path {
		t.Fatalf("Error.Path = %q; want %q", applyErr.Path, path)
	}
	if applyErr.Reason != reason {
		t.Fatalf("Error.Reason = %q; want %q", applyErr.Reason, reason)
	}
	if applyErr.Detail == "" {
		t.Fatalf("Error.Detail is empty")
	}
}

// requireStringMember checks a string object member.
func requireStringMember(t *testing.T, objectValue value.Value, name string, want string) {
	t.Helper()

	memberValue := requireMember(t, objectValue, name)
	got, ok := memberValue.String()
	if !ok {
		t.Fatalf("member %q kind = %s; want string", name, memberValue.Kind())
	}
	if got != want {
		t.Fatalf("member %q = %q; want %q", name, got, want)
	}
}

// requireMember returns an object member or fails the test.
func requireMember(t *testing.T, objectValue value.Value, name string) value.Value {
	t.Helper()

	view, ok := objectValue.Object()
	if !ok {
		t.Fatalf("value kind = %s; want object", objectValue.Kind())
	}

	memberValue, ok := view.Get(name)
	if !ok {
		t.Fatalf("member %q is absent", name)
	}

	return memberValue
}

// requireNoMember fails if an object member is present.
func requireNoMember(t *testing.T, objectValue value.Value, name string) {
	t.Helper()

	view, ok := objectValue.Object()
	if !ok {
		t.Fatalf("value kind = %s; want object", objectValue.Kind())
	}
	if view.Has(name) {
		t.Fatalf("member %q is present", name)
	}
}

// requireSet compares a fieldpath set by canonical string order.
func requireSet(t *testing.T, got fieldpath.Set, want ...string) {
	t.Helper()

	gotStrings := make([]string, 0, len(got.Paths()))
	for _, p := range got.Paths() {
		gotStrings = append(gotStrings, p.String())
	}

	if len(gotStrings) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(gotStrings, want) {
		t.Fatalf("set = %#v; want %#v", gotStrings, want)
	}
}

// requireOwners compares owners by deterministic string order.
func requireOwners(t *testing.T, got []fieldownership.Owner, want ...string) {
	t.Helper()

	gotStrings := make([]string, 0, len(got))
	for _, owner := range got {
		gotStrings = append(gotStrings, owner.String())
	}
	if len(gotStrings) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(gotStrings, want) {
		t.Fatalf("owners = %#v; want %#v", gotStrings, want)
	}
}

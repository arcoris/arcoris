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

package valueapply

import (
	"errors"
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func path(text string) fieldpath.Path {
	p, err := fieldpath.ParseCanonical(text)
	if err != nil {
		panic(err)
	}

	return p
}

func fields(paths ...fieldpath.Path) fieldpath.Set {
	return fieldpath.MustSet(paths...)
}

func root() fieldpath.Path {
	return fieldpath.Root()
}

func testFieldName(name string) fieldpath.FieldName {
	return fieldpath.MustFieldName(name)
}

func testMapKey(key string) fieldpath.MapKey {
	return fieldpath.MustMapKey(key)
}

func owner(name string) fieldownership.Owner {
	return fieldownership.MustOwner(name)
}

func entry(name string, paths ...fieldpath.Path) fieldownership.Entry {
	return fieldownership.MustEntry(owner(name), fields(paths...))
}

func state(entries ...fieldownership.Entry) fieldownership.State {
	return fieldownership.MustState(entries...)
}

func str(text string) value.Value {
	return value.StringValue(text)
}

func intValue(v int64) value.Value {
	return value.Int64Value(v)
}

func obj(members ...value.RecordMember) value.Value {
	return value.MustRecordValue(members...)
}

func list(items ...value.Value) value.Value {
	return value.MustListValue(items...)
}

func member(name string, v value.Value) value.RecordMember {
	return value.MustRecordMember(name, v)
}

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireSet(t *testing.T, got fieldpath.Set, want ...string) {
	t.Helper()

	gotStrings := pathStrings(got.Paths())
	if len(gotStrings) == 0 && len(want) == 0 {
		return
	}
	if !reflect.DeepEqual(gotStrings, want) {
		t.Fatalf("set = %#v; want %#v", gotStrings, want)
	}
}

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

func requireOwnersOf(t *testing.T, state fieldownership.State, path fieldpath.Path, want ...string) {
	t.Helper()

	got, err := state.OwnersOf(path)
	requireNoError(t, err)
	requireOwners(t, got, want...)
}

func requireStringMember(t *testing.T, object value.Value, name string, want string) {
	t.Helper()

	member := requireMember(t, object, name)
	got, ok := member.AsString()
	if !ok {
		t.Fatalf("member %q kind = %s; want string", name, member.Kind())
	}
	if got != want {
		t.Fatalf("member %q = %q; want %q", name, got, want)
	}
}

func requireNoMember(t *testing.T, object value.Value, name string) {
	t.Helper()

	view, ok := object.AsRecord()
	if !ok {
		t.Fatalf("value kind = %s; want record", object.Kind())
	}
	if view.Has(value.MemberName(name)) {
		t.Fatalf("member %q is present", name)
	}
}

func requireMember(t *testing.T, object value.Value, name string) value.Value {
	t.Helper()

	view, ok := object.AsRecord()
	if !ok {
		t.Fatalf("value kind = %s; want record", object.Kind())
	}

	member, ok := view.Get(value.MemberName(name))
	if !ok {
		t.Fatalf("member %q is absent", name)
	}

	return member
}

func requireListStrings(t *testing.T, got value.Value, want ...string) {
	t.Helper()

	view, ok := got.AsList()
	if !ok {
		t.Fatalf("value kind = %s; want list", got.Kind())
	}
	if view.Len() != len(want) {
		t.Fatalf("list length = %d; want %d", view.Len(), len(want))
	}

	for i, wantItem := range want {
		item, _ := view.At(i)
		gotItem, ok := item.AsString()
		if !ok {
			t.Fatalf("list[%d] kind = %s; want string", i, item.Kind())
		}
		if gotItem != wantItem {
			t.Fatalf("list[%d] = %q; want %q", i, gotItem, wantItem)
		}
	}
}

func requireListItem(t *testing.T, got value.Value, index int) value.Value {
	t.Helper()

	view, ok := got.AsList()
	if !ok {
		t.Fatalf("value kind = %s; want list", got.Kind())
	}

	item, ok := view.At(index)
	if !ok {
		t.Fatalf("list[%d] is absent", index)
	}

	return item
}

func pathStrings(paths []fieldpath.Path) []string {
	out := make([]string, 0, len(paths))
	for _, path := range paths {
		out = append(out, path.String())
	}

	return out
}

func specDescriptor() types.Descriptor {
	return types.Object(
		types.Field("image").String().Optional(),
		types.Field("replicas").Int64().Optional(),
	).Descriptor()
}

func specStringDescriptor() types.Descriptor {
	return types.Object(
		types.Field("image").String().Optional(),
		types.Field("replicas").String().Optional(),
	).Descriptor()
}

func typesUnknownPruneObject() types.Descriptor {
	return types.Object().UnknownFields(types.UnknownPrune).Descriptor()
}

func typesObjectWithSpec() types.Descriptor {
	return types.Object(
		types.Field("name").String().Optional(),
		types.Field("spec").Object(
			types.Field("image").String().Optional(),
		).Optional(),
	).Descriptor()
}

func mapDescriptor() types.Descriptor {
	return types.MapOf(types.String()).Descriptor()
}

func orderedStringListDescriptor() types.Descriptor {
	return types.ListOf(types.String()).Ordered().Descriptor()
}

func atomicStringListDescriptor() types.Descriptor {
	return types.ListOf(types.String()).Atomic().Descriptor()
}

func conditionsDescriptor() types.Descriptor {
	return types.ListOf(
		types.Object(
			types.Field("type").String().Required(),
			types.Field("status").String().Optional(),
		),
	).Map("type").Descriptor()
}

func specPath() fieldpath.Path {
	return path("$.spec")
}

func imagePath() fieldpath.Path {
	return path("$.image")
}

func replicasPath() fieldpath.Path {
	return path("$.replicas")
}

func labelPath() fieldpath.Path {
	return path(`$["app"]`)
}

func readySelector() fieldpath.Selector {
	return fieldpath.MustSelector(
		fieldpath.NewSelectorEntry(testFieldName("type"), fieldpath.StringLiteral("Ready")),
	)
}

func readyStatusPath() fieldpath.Path {
	return root().Select(readySelector()).Field(testFieldName("status"))
}

func readyCondition(status string) value.Value {
	return obj(member("type", str("Ready")), member("status", str(status)))
}

func specRequest(owner fieldownership.Owner) Request {
	return Request{
		Path:       root(),
		Owner:      owner,
		Live:       obj(member("image", str("api:v1")), member("replicas", str("3"))),
		Applied:    obj(member("image", str("api:v2"))),
		Descriptor: specStringDescriptor(),
		Ownership:  fieldownership.EmptyState(),
	}
}

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

package valuemerge

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

func requireErrorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}

func requireValue(t *testing.T, got value.Value, want value.Value) {
	t.Helper()

	if !valuesEqual(got, want) {
		t.Fatalf("value = %#v; want %#v", got, want)
	}
}

func requireStringMember(t *testing.T, object value.Value, name string, want string) {
	t.Helper()

	member := requireMember(t, object, name)
	text, ok := member.AsString()
	if !ok {
		t.Fatalf("member %q kind = %s; want string", name, member.Kind())
	}
	if text != want {
		t.Fatalf("member %q = %q; want %q", name, text, want)
	}
}

func requireIntegerMember(t *testing.T, object value.Value, name string, want int64) {
	t.Helper()

	member := requireMember(t, object, name)
	integer, ok := member.AsInteger()
	if !ok {
		t.Fatalf("member %q kind = %s; want integer", name, member.Kind())
	}

	n, ok := integer.Int64()
	if !ok {
		t.Fatalf("member %q does not fit int64", name)
	}
	if n != want {
		t.Fatalf("member %q = %d; want %d", name, n, want)
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

func requireRecordMemberOrder(t *testing.T, got value.Value, want ...string) {
	t.Helper()

	view, ok := got.AsRecord()
	if !ok {
		t.Fatalf("value kind = %s; want record", got.Kind())
	}
	if view.Len() != len(want) {
		t.Fatalf("record length = %d; want %d", view.Len(), len(want))
	}

	for i, member := range view.Members() {
		if member.Name.String() != want[i] {
			t.Fatalf("member[%d] = %q; want %q", i, member.Name, want[i])
		}
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

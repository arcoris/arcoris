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

package object

import (
	"testing"

	"arcoris.dev/apimachinery/api/meta"
)

func TestNewList(t *testing.T) {
	items := []int{1, 2}
	list := NewList(validListTypeMeta(), validListMeta(), items)

	if list.TypeMeta != validListTypeMeta() {
		t.Fatalf("TypeMeta = %#v", list.TypeMeta)
	}
	if list.ListMeta.ResourceVersion != "rv-1" {
		t.Fatalf("ListMeta = %#v", list.ListMeta)
	}
	if list.Len() != 2 {
		t.Fatalf("Len() = %d", list.Len())
	}

	items[0] = 9
	if list.Items[0] != 1 {
		t.Fatalf("NewList() did not detach item slice")
	}
}

func TestNewListPreservesNilAndEmptyItems(t *testing.T) {
	nilList := NewList[int](validListTypeMeta(), validListMeta(), nil)
	if nilList.Items != nil {
		t.Fatalf("nil items became %#v", nilList.Items)
	}

	emptyList := NewList(validListTypeMeta(), validListMeta(), []int{})
	if emptyList.Items == nil || len(emptyList.Items) != 0 {
		t.Fatalf("empty items became %#v", emptyList.Items)
	}
}

func TestNewListItemCopyIsShallow(t *testing.T) {
	type item struct {
		Values []int
	}

	items := []item{{Values: []int{1}}}
	list := NewList(validListTypeMeta(), validListMeta(), items)
	items[0].Values[0] = 9

	if list.Items[0].Values[0] != 9 {
		t.Fatal("NewList() unexpectedly deep-copied item values")
	}
}

func TestListIsZero(t *testing.T) {
	tests := []struct {
		name string
		list List[int]
		want bool
	}{
		{name: "zero", list: List[int]{}, want: true},
		{name: "type meta", list: List[int]{TypeMeta: validListTypeMeta()}, want: false},
		{name: "list meta", list: List[int]{ListMeta: meta.ListMeta{ResourceVersion: "rv-1"}}, want: false},
		{name: "items", list: List[int]{Items: []int{0}}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.IsZero(); got != tt.want {
				t.Fatalf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

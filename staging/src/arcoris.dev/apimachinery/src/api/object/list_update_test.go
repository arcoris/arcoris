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

func TestListUpdateMethods(t *testing.T) {
	original := NewList(validListTypeMeta(), validPageMeta(), []int{1, 2})

	otherType := meta.TypeMeta{Kind: "OtherList"}
	withType := original.WithTypeMeta(otherType)
	if withType.TypeMeta != otherType {
		t.Fatalf("WithTypeMeta() TypeMeta = %#v", withType.TypeMeta)
	}
	if original.TypeMeta != validListTypeMeta() {
		t.Fatal("WithTypeMeta() mutated original")
	}

	count := uint64(2)
	otherPage := meta.PageMeta{
		ResourceVersion:    "rv-2",
		RemainingItemCount: &count,
	}
	withPage := original.WithPageMeta(otherPage)
	if withPage.PageMeta.ResourceVersion != "rv-2" {
		t.Fatalf("WithPageMeta() PageMeta = %#v", withPage.PageMeta)
	}
	if original.PageMeta.ResourceVersion != "rv-1" {
		t.Fatal("WithPageMeta() mutated original")
	}
	*otherPage.RemainingItemCount = 9
	if *withPage.PageMeta.RemainingItemCount != 2 {
		t.Fatal("WithPageMeta() did not detach replacement metadata")
	}

	items := []int{3, 4}
	withItems := original.WithItems(items)
	if len(withItems.Items) != 2 || withItems.Items[0] != 3 || withItems.Items[1] != 4 {
		t.Fatalf("WithItems() Items = %#v", withItems.Items)
	}
	if len(original.Items) != 2 || original.Items[0] != 1 || original.Items[1] != 2 {
		t.Fatal("WithItems() mutated original")
	}

	items[0] = 9
	if withItems.Items[0] != 3 {
		t.Fatal("WithItems() did not detach caller slice")
	}
}

func TestListWithItemsPreservesNilAndEmpty(t *testing.T) {
	original := NewList(validListTypeMeta(), validPageMeta(), []int{1})

	nilItems := original.WithItems(nil)
	if nilItems.Items != nil {
		t.Fatalf("WithItems(nil) Items = %#v", nilItems.Items)
	}

	emptyItems := original.WithItems([]int{})
	if emptyItems.Items == nil || len(emptyItems.Items) != 0 {
		t.Fatalf("WithItems(empty) Items = %#v", emptyItems.Items)
	}
}

func TestListWithItemsIsShallow(t *testing.T) {
	type item struct {
		Values []int
	}

	items := []item{{Values: []int{1}}}
	list := NewList(validListTypeMeta(), validPageMeta(), []item{})
	withItems := list.WithItems(items)
	items[0].Values[0] = 9

	if withItems.Items[0].Values[0] != 9 {
		t.Fatal("WithItems() unexpectedly deep-copied item values")
	}
}

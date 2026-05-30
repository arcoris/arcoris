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

import "testing"

func TestListItemsCopy(t *testing.T) {
	if got := (List[int]{}).ItemsCopy(); got != nil {
		t.Fatalf("nil ItemsCopy() = %#v", got)
	}

	empty := List[int]{Items: []int{}}
	emptyCopy := empty.ItemsCopy()
	if emptyCopy == nil || len(emptyCopy) != 0 {
		t.Fatalf("empty ItemsCopy() = %#v", emptyCopy)
	}

	list := List[int]{Items: []int{1, 2}}
	items := list.ItemsCopy()
	if len(items) != 2 || items[0] != 1 || items[1] != 2 {
		t.Fatalf("ItemsCopy() = %#v", items)
	}

	items[0] = 9
	if list.Items[0] != 1 {
		t.Fatal("ItemsCopy() mutated original slice")
	}
}

func TestListItemsCopyIsShallow(t *testing.T) {
	type item struct {
		Values []int
	}

	list := List[item]{Items: []item{{Values: []int{1}}}}
	items := list.ItemsCopy()
	items[0].Values[0] = 9

	if list.Items[0].Values[0] != 9 {
		t.Fatal("ItemsCopy() unexpectedly deep-copied item values")
	}
}

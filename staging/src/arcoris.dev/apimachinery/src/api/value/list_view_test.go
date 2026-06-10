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

package value

import "testing"

func TestListViewAccessors(t *testing.T) {
	value := mustList(t, StringValue("first"), BytesValue([]byte{1, 2}))
	view, ok := value.AsList()
	requireEqual(t, ok, true)

	requireEqual(t, view.Len(), 2)
	requireEqual(t, view.IsEmpty(), false)

	first, ok := view.At(0)
	requireEqual(t, ok, true)

	text, ok := first.AsString()
	requireEqual(t, ok, true)
	requireEqual(t, text, "first")

	_, ok = view.At(-1)
	requireEqual(t, ok, false)

	_, ok = view.At(2)
	requireEqual(t, ok, false)
}

func TestListViewEmptyItemsIsNonNil(t *testing.T) {
	value := mustList(t)
	view, ok := value.AsList()
	requireEqual(t, ok, true)

	items := view.Items()
	if items == nil {
		t.Fatal("Items() returned nil")
	}
	requireEqual(t, len(items), 0)
}

func TestListViewItemsReturnsDetachedSliceWithoutDeepCloningValues(t *testing.T) {
	value := mustList(t, BytesValue([]byte{1, 2}))
	view, ok := value.AsList()
	requireEqual(t, ok, true)

	items := view.Items()
	items[0].bytesValue[0] = 9

	item, ok := view.At(0)
	requireEqual(t, ok, true)

	bytes, ok := item.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{9, 2})
}

func TestListViewCloneItemsDeepClonesValues(t *testing.T) {
	value := mustList(t, BytesValue([]byte{1, 2}))
	view, ok := value.AsList()
	requireEqual(t, ok, true)

	items := view.CloneItems()
	items[0].bytesValue[0] = 9

	item, ok := view.At(0)
	requireEqual(t, ok, true)
	bytes, ok := item.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestListViewCloneAtDeepClonesValue(t *testing.T) {
	value := mustList(t, BytesValue([]byte{1, 2}))
	view, ok := value.AsList()
	requireEqual(t, ok, true)

	item, ok := view.CloneAt(0)
	requireEqual(t, ok, true)
	item.bytesValue[0] = 9

	again, ok := view.At(0)
	requireEqual(t, ok, true)
	bytes, ok := again.AsBytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestListViewForEachVisitsInOrderAndStopsEarly(t *testing.T) {
	value := mustList(t, StringValue("first"), StringValue("second"))
	view, ok := value.AsList()
	requireEqual(t, ok, true)

	var texts []string
	view.ForEach(func(index int, item Value) bool {
		text, ok := item.AsString()
		requireEqual(t, ok, true)
		texts = append(texts, text)
		return false
	})

	requireStringsEqual(t, texts, []string{"first"})
}

func TestListWrongKindAccessorReturnsFalse(t *testing.T) {
	_, ok := NullValue().AsList()
	requireEqual(t, ok, false)
}

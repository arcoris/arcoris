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

func TestMapViewAccessors(t *testing.T) {
	value := mustMap(t,
		MapEntry("name", String("worker")),
		MapEntry("payload", Bytes([]byte{1, 2})),
	)
	view, ok := value.Map()
	requireEqual(t, ok, true)

	requireEqual(t, view.Len(), 2)
	requireEqual(t, view.IsEmpty(), false)
	requireEqual(t, view.Has("name"), true)
	requireEqual(t, view.Has("missing"), false)
	requireStringsEqual(t, view.Keys(), []string{"name", "payload"})

	got, ok := view.Get("name")
	requireEqual(t, ok, true)

	name, ok := got.String()
	requireEqual(t, ok, true)
	requireEqual(t, name, "worker")

	_, ok = view.Get("missing")
	requireEqual(t, ok, false)
}

func TestMapViewReturnsDetachedEntriesAndValues(t *testing.T) {
	value := mustMap(t, MapEntry("payload", Bytes([]byte{1, 2})))
	view, ok := value.Map()
	requireEqual(t, ok, true)

	entries := view.Entries()
	entries[0].Key = "changed"
	entries[0].Value.bytesValue[0] = 9

	requireEqual(t, view.Has("payload"), true)
	got, ok := view.Get("payload")
	requireEqual(t, ok, true)

	bytes, ok := got.Bytes()
	requireEqual(t, ok, true)
	requireBytesEqual(t, bytes, []byte{1, 2})
}

func TestMapWrongKindAccessorReturnsFalse(t *testing.T) {
	_, ok := Null().Map()
	requireEqual(t, ok, false)
}

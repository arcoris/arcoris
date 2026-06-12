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

package objectownership

import (
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/fieldownership"
)

func TestSurfaceIsEmpty(t *testing.T) {
	if !(Surface{}).IsEmpty() {
		t.Fatalf("empty surface IsEmpty() = false")
	}
	if !(Surface{Entries: []Entry{documentEntry("user")}}).IsEmpty() {
		t.Fatalf("surface with empty entry IsEmpty() = false")
	}
	if (Surface{Entries: []Entry{documentEntry("user", "$.image")}}).IsEmpty() {
		t.Fatalf("surface with owned field IsEmpty() = true")
	}
}

func TestSurfaceIsEmptyDoesNotValidateEntries(t *testing.T) {
	surface := Surface{
		Entries: []Entry{
			{Owner: fieldownership.Owner{}, Fields: []Path{""}},
		},
	}

	if surface.IsEmpty() {
		t.Fatalf("surface with invalid but mentioned field IsEmpty() = true")
	}
}

func TestSurfaceClonePreservesRawShape(t *testing.T) {
	surface := Surface{Entries: []Entry{
		documentEntry("b", "$.b", "$.b"),
		documentEntry("a"),
		documentEntry("b", "$.a"),
	}}

	clone := surface.Clone()

	if !reflect.DeepEqual(clone, surface) {
		t.Fatalf("Clone() = %#v; want %#v", clone, surface)
	}
}

func TestSurfaceClonePreservesNilAndEmptyEntries(t *testing.T) {
	nilEntries := (Surface{}).Clone()
	if nilEntries.Entries != nil {
		t.Fatalf("nil Entries clone = %#v; want nil", nilEntries.Entries)
	}

	emptyEntries := (Surface{Entries: []Entry{}}).Clone()
	if emptyEntries.Entries == nil {
		t.Fatal("empty Entries clone = nil; want non-nil empty slice")
	}
}

func TestSurfaceCloneDetachesEntriesAndFields(t *testing.T) {
	surface := Surface{Entries: []Entry{documentEntry("user", "$.image")}}

	clone := surface.Clone()
	clone.Entries[0].Owner = owner("other")
	clone.Entries[0].Fields[0] = "$.other"

	requireDocumentEntries(t, surface, documentEntry("user", "$.image"))
}

func TestSurfaceEntriesCopyDetachesEntriesAndFields(t *testing.T) {
	surface := Surface{Entries: []Entry{documentEntry("user", "$.image")}}

	entries := surface.EntriesCopy()
	entries[0].Owner = owner("other")
	entries[0].Fields[0] = "$.other"

	requireDocumentEntries(t, surface, documentEntry("user", "$.image"))
}

func TestSurfaceEntriesCopyPreservesNilAndEmptyEntries(t *testing.T) {
	if entries := (Surface{}).EntriesCopy(); entries != nil {
		t.Fatalf("nil EntriesCopy() = %#v; want nil", entries)
	}

	entries := (Surface{Entries: []Entry{}}).EntriesCopy()
	if entries == nil {
		t.Fatal("empty EntriesCopy() = nil; want non-nil empty slice")
	}
	if len(entries) != 0 {
		t.Fatalf("empty EntriesCopy() len = %d; want 0", len(entries))
	}
}

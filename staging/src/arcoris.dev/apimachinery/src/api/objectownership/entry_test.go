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

func TestEntryIsEmpty(t *testing.T) {
	if !documentEntry("user").IsEmpty() {
		t.Fatalf("empty entry IsEmpty() = false")
	}
	if documentEntry("user", "$.image").IsEmpty() {
		t.Fatalf("entry with field IsEmpty() = true")
	}
}

func TestEntryIsEmptyDoesNotValidateOwner(t *testing.T) {
	entry := Entry{Owner: fieldownership.Owner{}}

	if !entry.IsEmpty() {
		t.Fatalf("entry with invalid owner but no fields IsEmpty() = false")
	}
}

func TestEntryClonePreservesRawFields(t *testing.T) {
	entry := documentEntry("user", "$.z", "$.z", "$.a")

	clone := entry.Clone()

	if !reflect.DeepEqual(clone, entry) {
		t.Fatalf("Clone() = %#v; want %#v", clone, entry)
	}
}

func TestEntryClonePreservesNilAndEmptyFields(t *testing.T) {
	nilFields := (Entry{Owner: owner("user")}).Clone()
	if nilFields.Fields != nil {
		t.Fatalf("nil Fields clone = %#v; want nil", nilFields.Fields)
	}

	emptyFields := (Entry{Owner: owner("user"), Fields: []Path{}}).Clone()
	if emptyFields.Fields == nil {
		t.Fatal("empty Fields clone = nil; want non-nil empty slice")
	}
}

func TestEntryCloneDetachesFields(t *testing.T) {
	entry := documentEntry("user", "$.image")

	clone := entry.Clone()
	clone.Fields[0] = "$.other"

	if entry.Fields[0] != "$.image" {
		t.Fatalf("original field = %q; want $.image", entry.Fields[0])
	}
}

func TestEntryFieldsCopyDetachesFields(t *testing.T) {
	entry := documentEntry("user", "$.image")

	fields := entry.FieldsCopy()
	fields[0] = "$.other"

	if entry.Fields[0] != "$.image" {
		t.Fatalf("original field = %q; want $.image", entry.Fields[0])
	}
}

func TestEntryFieldsCopyPreservesNilAndEmptyFields(t *testing.T) {
	if fields := (Entry{Owner: owner("user")}).FieldsCopy(); fields != nil {
		t.Fatalf("nil FieldsCopy() = %#v; want nil", fields)
	}

	fields := (Entry{Owner: owner("user"), Fields: []Path{}}).FieldsCopy()
	if fields == nil {
		t.Fatal("empty FieldsCopy() = nil; want non-nil empty slice")
	}
	if len(fields) != 0 {
		t.Fatalf("empty FieldsCopy() len = %d; want 0", len(fields))
	}
}

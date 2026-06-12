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
)

func TestDocumentIsEmpty(t *testing.T) {
	if !(Document{Version: DocumentVersionV1}).IsEmpty() {
		t.Fatalf("empty document IsEmpty() = false")
	}
	if !document(documentEntry("user")).IsEmpty() {
		t.Fatalf("document with empty entry IsEmpty() = false")
	}
	if document(documentEntry("user", "$.image")).IsEmpty() {
		t.Fatalf("document with entry IsEmpty() = true")
	}
}

func TestDocumentIsEmptyIgnoresVersionValidity(t *testing.T) {
	if !(Document{Version: "unsupported"}).IsEmpty() {
		t.Fatalf("document with invalid version but no ownership IsEmpty() = false")
	}
}

func TestDocumentCloneZeroDocument(t *testing.T) {
	got := (Document{}).Clone()

	if !reflect.DeepEqual(got, Document{}) {
		t.Fatalf("Clone() = %#v; want zero document", got)
	}
}

func TestDocumentClonePreservesRawShape(t *testing.T) {
	doc := Document{
		Version: "raw-version",
		Desired: Surface{Entries: []Entry{
			documentEntry("b", "$.b", "$.b"),
			documentEntry("a"),
			documentEntry("b", "$.a"),
		}},
	}

	clone := doc.Clone()

	if !reflect.DeepEqual(clone, doc) {
		t.Fatalf("Clone() = %#v; want %#v", clone, doc)
	}
}

func TestDocumentClonePreservesNilAndEmptyEntries(t *testing.T) {
	nilEntries := Document{Version: DocumentVersionV1}.Clone()
	if nilEntries.Desired.Entries != nil {
		t.Fatalf("nil Entries clone = %#v; want nil", nilEntries.Desired.Entries)
	}

	emptyEntries := Document{
		Version: DocumentVersionV1,
		Desired: Surface{Entries: []Entry{}},
	}.Clone()
	if emptyEntries.Desired.Entries == nil {
		t.Fatal("empty Entries clone = nil; want non-nil empty slice")
	}
	if len(emptyEntries.Desired.Entries) != 0 {
		t.Fatalf("empty Entries len = %d; want 0", len(emptyEntries.Desired.Entries))
	}
}

func TestDocumentCloneDetachesEntriesAndFields(t *testing.T) {
	doc := document(documentEntry("user", "$.image"))

	clone := doc.Clone()
	clone.Desired.Entries[0].Owner = owner("other")
	clone.Desired.Entries[0].Fields[0] = "$.other"

	requireDocumentEntries(t, doc.Desired, documentEntry("user", "$.image"))

	doc.Desired.Entries[0].Owner = owner("source")
	doc.Desired.Entries[0].Fields[0] = "$.source"

	requireDocumentEntries(t, clone.Desired, documentEntry("other", "$.other"))
}

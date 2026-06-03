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

import "testing"

func TestDocumentIsEmpty(t *testing.T) {
	if !(Document{Version: VersionV1}).IsEmpty() {
		t.Fatalf("empty document IsEmpty() = false")
	}
	if !document(documentEntry("user")).IsEmpty() {
		t.Fatalf("document with empty entry IsEmpty() = false")
	}
	if document(documentEntry("user", "$.image")).IsEmpty() {
		t.Fatalf("document with entry IsEmpty() = true")
	}
}

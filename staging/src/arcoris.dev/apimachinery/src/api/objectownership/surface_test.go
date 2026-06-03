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
			documentEntry(" ", ""),
		},
	}

	if surface.IsEmpty() {
		t.Fatalf("surface with invalid but mentioned field IsEmpty() = true")
	}
}

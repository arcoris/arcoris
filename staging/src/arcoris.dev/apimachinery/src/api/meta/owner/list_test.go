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

package owner

import "testing"

func TestList(t *testing.T) {
	var nilList List
	if !nilList.IsZero() {
		t.Fatal("nil List IsZero() = false")
	}
	if nilList.Len() != 0 {
		t.Fatalf("nil List Len() = %d, want 0", nilList.Len())
	}

	owners := List{validReference(true)}
	if owners.IsZero() {
		t.Fatal("non-zero List IsZero() = true")
	}
	if owners.Len() != 1 {
		t.Fatalf("Len() = %d", owners.Len())
	}
}

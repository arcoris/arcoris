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

package labels

import "testing"

func TestSet(t *testing.T) {
	var nilSet Set
	if !nilSet.IsZero() {
		t.Fatal("nil Set IsZero() = false")
	}
	if nilSet.Len() != 0 {
		t.Fatalf("nil Set Len() = %d, want 0", nilSet.Len())
	}

	set := Set{"role": "worker"}
	if set.IsZero() {
		t.Fatal("non-zero Set IsZero() = true")
	}
	if set.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", set.Len())
	}
	if !set.Has("role") {
		t.Fatal("Has() = false")
	}
	if value, ok := set.Get("role"); !ok || value != "worker" {
		t.Fatalf("Get() = %q, %v", value, ok)
	}
}

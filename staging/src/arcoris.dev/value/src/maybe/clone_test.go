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

package maybe_test

import (
	"slices"
	"testing"

	"arcoris.dev/value/maybe"
)

func cloneStrings(val []string) []string {
	return slices.Clone(val)
}

func TestCloneSome(t *testing.T) {
	original := []string{"a", "b"}
	m := maybe.Some(original)

	cloned := m.Clone(cloneStrings)
	got, ok := cloned.Load()
	if !ok {
		t.Fatal("Clone returned None for Some")
	}
	if got[0] != "a" || got[1] != "b" {
		t.Fatalf("cloned value = %#v, want %#v", got, original)
	}

	original[0] = "changed"
	if got[0] != "a" {
		t.Fatalf("Clone did not isolate value: got[0] = %q", got[0])
	}
}

func TestCloneNone(t *testing.T) {
	cloned := maybe.None[[]string]().Clone(cloneStrings)
	if !cloned.IsNone() {
		t.Fatal("Clone returned Some for None")
	}
}

func TestClonePanicsOnNilClone(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Clone did not panic for nil clone function")
		}
	}()

	_ = maybe.None[string]().Clone(nil)
}

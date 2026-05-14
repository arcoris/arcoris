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

package result_test

import (
	"errors"
	"slices"
	"testing"

	"arcoris.dev/value/result"
)

func cloneStrings(value []string) []string {
	return slices.Clone(value)
}

func TestCloneOK(t *testing.T) {
	original := []string{"a", "b"}
	r := result.OK(original)

	cloned := r.Clone(cloneStrings)
	got, err := cloned.Load()
	if err != nil {
		t.Fatalf("Clone returned Err for OK: %v", err)
	}
	if got[0] != "a" || got[1] != "b" {
		t.Fatalf("cloned value = %#v, want %#v", got, original)
	}

	original[0] = "changed"
	if got[0] != "a" {
		t.Fatalf("Clone did not isolate value: got[0] = %q", got[0])
	}
}

func TestCloneErr(t *testing.T) {
	want := errors.New("failed")
	cloned := result.Err[[]string](want).Clone(cloneStrings)

	_, err := cloned.Load()
	if !errors.Is(err, want) {
		t.Fatalf("Clone error = %v, want %v", err, want)
	}
}

func TestClonePanicsOnNilClone(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Clone did not panic for nil clone function")
		}
	}()

	_ = result.Err[string](errors.New("failed")).Clone(nil)
}

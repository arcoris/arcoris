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


package merge

import "testing"

func TestLinearMergesInOrder(t *testing.T) {
	got, ok := Linear([]string{"a", "b", "c", "d"}, func(dst *string, src string) {
		*dst = "(" + *dst + "+" + src + ")"
	})
	if !ok {
		t.Fatal("expected non-empty merge")
	}
	if got != "(((a+b)+c)+d)" {
		t.Fatalf("got %q, want linear left fold", got)
	}
}

func TestLinearDoesNotMutateInput(t *testing.T) {
	partials := []string{"a", "b", "c"}
	_, ok := Linear(partials, func(dst *string, src string) { *dst += src })
	if !ok {
		t.Fatal("expected non-empty merge")
	}
	want := []string{"a", "b", "c"}
	for i := range want {
		if partials[i] != want[i] {
			t.Fatalf("partials = %#v, want %#v", partials, want)
		}
	}
}

func TestLinearEmpty(t *testing.T) {
	_, ok := Linear[int](nil, func(dst *int, src int) { *dst += src })
	if ok {
		t.Fatal("empty merge returned ok")
	}
}

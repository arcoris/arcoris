/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package merge

import (
	"testing"

	"arcoris.dev/measure/internal/reduce"
)

func TestMergeDispatchesLinear(t *testing.T) {
	got, ok := Merge([]string{"a", "b"}, reduce.MergeLinear, func(dst *string, src string) { *dst += src })
	if !ok || got != "ab" {
		t.Fatalf("got %q ok=%v, want ab true", got, ok)
	}
}

func TestMergeDispatchesPairwise(t *testing.T) {
	got, ok := Merge([]int{1, 2, 3, 4}, reduce.MergePairwise, func(dst *int, src int) { *dst += src })
	if !ok || got != 10 {
		t.Fatalf("got %d ok=%v, want 10 true", got, ok)
	}
}

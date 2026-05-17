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

package core

import "testing"

func TestMapperCallbacks(t *testing.T) {
	mapper := Mapper[int](func(r Range) int { return r.Len() })
	into := IntoMapper[int](func(r Range, dst *int) { *dst += r.Len() })
	indexed := IndexedIntoMapper[int](func(worker int, r Range, dst *int) {
		*dst += worker + r.Len()
	})

	got := mapper(Range{Start: 0, End: 3})
	into(Range{Start: 0, End: 2}, &got)
	indexed(4, Range{Start: 0, End: 1}, &got)
	if got != 10 {
		t.Fatalf("callback result = %d, want 10", got)
	}
}

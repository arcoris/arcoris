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

package reduce

import "testing"

func TestScratchBuffers(t *testing.T) {
	s := Scratch[int]{
		Ranges:   make([]Range, 3, 8),
		Partials: []int{1, 2, 3},
	}
	rangeCap := cap(s.Ranges)
	partialCap := cap(s.Partials)

	s.Reset()
	if len(s.Ranges) != 0 || len(s.Partials) != 0 {
		t.Fatalf("Reset lengths = ranges:%d partials:%d, want zero", len(s.Ranges), len(s.Partials))
	}
	if cap(s.Ranges) != rangeCap || cap(s.Partials) != partialCap {
		t.Fatal("Reset changed backing capacity")
	}

	s.Partials = append(s.Partials, 7, 8, 9)
	zeroed := s.EnsurePartials(2)
	if len(zeroed) != 2 || zeroed[0] != 0 || zeroed[1] != 0 {
		t.Fatalf("EnsurePartials() = %#v, want two zeroed slots", zeroed)
	}

	s.Partials = append(s.Partials[:0], 7, 8, 9)
	dirty := s.EnsurePartialsDirty(2)
	if len(dirty) != 2 || dirty[0] != 7 || dirty[1] != 8 {
		t.Fatalf("EnsurePartialsDirty() = %#v, want reused dirty slots", dirty)
	}
}

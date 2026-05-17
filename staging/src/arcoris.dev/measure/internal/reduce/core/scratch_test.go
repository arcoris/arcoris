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

func TestScratchResetKeepsBackingStorage(t *testing.T) {
	s := Scratch[int]{
		Ranges:   make([]Range, 3, 8),
		Partials: []int{1, 2, 3},
		Used:     []bool{true, false, true},
	}
	rangeCap := cap(s.Ranges)
	partialCap := cap(s.Partials)
	usedCap := cap(s.Used)

	s.Reset()
	if len(s.Ranges) != 0 || len(s.Partials) != 0 || len(s.Used) != 0 {
		t.Fatalf("Reset lengths = ranges:%d partials:%d used:%d, want zero", len(s.Ranges), len(s.Partials), len(s.Used))
	}
	if cap(s.Ranges) != rangeCap || cap(s.Partials) != partialCap || cap(s.Used) != usedCap {
		t.Fatal("Reset changed backing capacity")
	}
}

func TestScratchResetCanRetainPointers(t *testing.T) {
	value := 42
	s := Scratch[*int]{Partials: []*int{&value}}
	s.Reset()

	retained := s.Partials[:cap(s.Partials)]
	if retained[0] != &value {
		t.Fatal("Reset cleared backing pointer; want retained backing storage")
	}
}

func TestScratchClearZeroesVisibleSlots(t *testing.T) {
	value := 42
	s := Scratch[*int]{
		Ranges:   []Range{{Start: 1, End: 3}},
		Partials: []*int{&value},
		Used:     []bool{true},
	}
	ranges := s.Ranges
	partials := s.Partials
	used := s.Used
	rangeCap := cap(s.Ranges)
	partialCap := cap(s.Partials)
	usedCap := cap(s.Used)

	s.Clear()
	if len(s.Ranges) != 0 || len(s.Partials) != 0 || len(s.Used) != 0 {
		t.Fatalf("Clear lengths = ranges:%d partials:%d used:%d, want zero", len(s.Ranges), len(s.Partials), len(s.Used))
	}
	if cap(s.Ranges) != rangeCap || cap(s.Partials) != partialCap || cap(s.Used) != usedCap {
		t.Fatal("Clear changed backing capacity")
	}
	if ranges[0] != (Range{}) {
		t.Fatalf("range backing slot = %#v, want zero range", ranges[0])
	}
	if partials[0] != nil {
		t.Fatal("partial backing slot retained pointer after Clear")
	}
	if used[0] {
		t.Fatal("used backing slot stayed true after Clear")
	}
}

func TestScratchClearZeroesPointersAfterReset(t *testing.T) {
	value := 42
	s := Scratch[*int]{Partials: []*int{&value}}
	s.Reset()
	s.Clear()

	retained := s.Partials[:cap(s.Partials)]
	if retained[0] != nil {
		t.Fatal("Clear retained pointer hidden by Reset")
	}
}

func TestScratchReleaseDropsBackingStorage(t *testing.T) {
	value := 42
	s := Scratch[*int]{
		Ranges:   []Range{{Start: 1, End: 2}},
		Partials: []*int{&value},
		Used:     []bool{true},
	}

	s.Release()
	if s.Ranges != nil || s.Partials != nil || s.Used != nil {
		t.Fatalf("Release slices = ranges:%#v partials:%#v used:%#v, want nil", s.Ranges, s.Partials, s.Used)
	}
}

func TestScratchEnsurePartialsZeroesSlots(t *testing.T) {
	s := Scratch[int]{Partials: []int{7, 8, 9}}
	zeroed := s.EnsurePartials(2)
	if len(zeroed) != 2 || zeroed[0] != 0 || zeroed[1] != 0 {
		t.Fatalf("EnsurePartials() = %#v, want two zeroed slots", zeroed)
	}
}

func TestScratchEnsurePartialsDirtyKeepsSlots(t *testing.T) {
	s := Scratch[int]{Partials: []int{7, 8, 9}}
	dirty := s.EnsurePartialsDirty(2)
	if len(dirty) != 2 || dirty[0] != 7 || dirty[1] != 8 {
		t.Fatalf("EnsurePartialsDirty() = %#v, want reused dirty slots", dirty)
	}
}

func TestScratchResetKeepsUsedBackingStorage(t *testing.T) {
	s := Scratch[int]{Used: []bool{true, false, true}}
	usedCap := cap(s.Used)

	s.Reset()
	if len(s.Used) != 0 {
		t.Fatalf("Used len = %d, want 0", len(s.Used))
	}
	if cap(s.Used) != usedCap {
		t.Fatal("Reset changed Used capacity")
	}
}

func TestScratchClearZeroesUsedSlots(t *testing.T) {
	s := Scratch[int]{Used: []bool{true, false, true}}
	used := s.Used

	s.Clear()
	for i, got := range used {
		if got {
			t.Fatalf("used[%d] = true, want false", i)
		}
	}
}

func TestScratchReleaseDropsUsedStorage(t *testing.T) {
	s := Scratch[int]{Used: []bool{true}}

	s.Release()
	if s.Used != nil {
		t.Fatalf("Used = %#v, want nil", s.Used)
	}
}

func TestScratchEnsureUsedZeroesSlots(t *testing.T) {
	s := Scratch[int]{Used: []bool{true, true, false}}

	got := s.EnsureUsed(2)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	for i, used := range got {
		if used {
			t.Fatalf("used[%d] = true, want false", i)
		}
	}
}

func TestScratchEnsureUsedDirtyKeepsSlots(t *testing.T) {
	s := Scratch[int]{Used: []bool{true, false, true}}

	got := s.EnsureUsedDirty(2)
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if !got[0] || got[1] {
		t.Fatalf("EnsureUsedDirty() = %#v, want dirty reused slots", got)
	}
}

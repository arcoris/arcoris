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

package atomicx

import (
	"sync/atomic"
	"testing"
	"unsafe"
)

type padTestValue struct {
	value int
}

// TestCacheLinePadSize verifies the package's fixed false-sharing pad width.
func TestCacheLinePadSize(t *testing.T) {
	t.Parallel()

	if CacheLinePadSize != 64 {
		t.Fatalf("CacheLinePadSize = %d, want 64", CacheLinePadSize)
	}
}

// TestCacheLinePadMemorySize verifies CacheLinePad occupies exactly the
// configured width, not merely a type-level constant with a matching value.
func TestCacheLinePadMemorySize(t *testing.T) {
	t.Parallel()

	got := unsafe.Sizeof(CacheLinePad{})
	want := uintptr(CacheLinePadSize)

	if got != want {
		t.Fatalf("unsafe.Sizeof(CacheLinePad{}) = %d, want %d", got, want)
	}
}

// TestAtomicWrapperValueSizes verifies production layout constants against
// sync/atomic wrapper sizes.
func TestAtomicWrapperValueSizes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		got  uintptr
		want uintptr
	}{
		{
			name: "atomic.Uint64",
			got:  unsafe.Sizeof(atomic.Uint64{}),
			want: uintptr(atomicUint64Size),
		},
		{
			name: "atomic.Uint32",
			got:  unsafe.Sizeof(atomic.Uint32{}),
			want: uintptr(atomicUint32Size),
		},
		{
			name: "atomic.Int64",
			got:  unsafe.Sizeof(atomic.Int64{}),
			want: uintptr(atomicInt64Size),
		},
		{
			name: "atomic.Int32",
			got:  unsafe.Sizeof(atomic.Int32{}),
			want: uintptr(atomicInt32Size),
		},
		{
			name: "atomic.Pointer",
			got:  unsafe.Sizeof(atomic.Pointer[padTestValue]{}),
			want: uintptr(atomicPointerSize),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.got != tc.want {
				t.Fatalf("unsafe.Sizeof(%s{}) = %d, want %d", tc.name, tc.got, tc.want)
			}
		})
	}
}

// TestPaddedPrimitiveSizesIncludeLeadingPadAndValueSlot verifies primitive
// layouts have enough total space for a leading pad and completed value slot.
func TestPaddedPrimitiveSizesIncludeLeadingPadAndValueSlot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size uintptr
		min  uintptr
	}{
		{
			name: "PaddedUint64",
			size: unsafe.Sizeof(PaddedUint64{}),
			min: uintptr(CacheLinePadSize + atomicUint64Size +
				(CacheLinePadSize - atomicUint64Size)),
		},
		{
			name: "PaddedUint32",
			size: unsafe.Sizeof(PaddedUint32{}),
			min: uintptr(CacheLinePadSize + atomicUint32Size +
				(CacheLinePadSize - atomicUint32Size)),
		},
		{
			name: "PaddedInt64",
			size: unsafe.Sizeof(PaddedInt64{}),
			min: uintptr(CacheLinePadSize + atomicInt64Size +
				(CacheLinePadSize - atomicInt64Size)),
		},
		{
			name: "PaddedInt32",
			size: unsafe.Sizeof(PaddedInt32{}),
			min: uintptr(CacheLinePadSize + atomicInt32Size +
				(CacheLinePadSize - atomicInt32Size)),
		},
		{
			name: "PaddedPointer",
			size: unsafe.Sizeof(PaddedPointer[padTestValue]{}),
			min: uintptr(CacheLinePadSize + atomicPointerSize +
				(CacheLinePadSize - atomicPointerSize)),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.size < tc.min {
				t.Fatalf("unsafe.Sizeof(%s{}) = %d, want at least %d", tc.name, tc.size, tc.min)
			}
		})
	}
}

// TestPaddedPrimitiveValueOffsetsKeepLeadingPad verifies each atomic wrapper is
// placed after the leading pad. Total struct size alone would not prove that the
// hot value is isolated from the previous field.
func TestPaddedPrimitiveValueOffsetsKeepLeadingPad(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		offset uintptr
	}{
		{name: "PaddedUint64", offset: unsafe.Offsetof(PaddedUint64{}.value)},
		{name: "PaddedUint32", offset: unsafe.Offsetof(PaddedUint32{}.value)},
		{name: "PaddedInt64", offset: unsafe.Offsetof(PaddedInt64{}.value)},
		{name: "PaddedInt32", offset: unsafe.Offsetof(PaddedInt32{}.value)},
		{name: "PaddedPointer", offset: unsafe.Offsetof(PaddedPointer[padTestValue]{}.value)},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.offset < uintptr(CacheLinePadSize) {
				t.Fatalf("%s.value offset = %d, want at least %d", tc.name, tc.offset, CacheLinePadSize)
			}
		})
	}
}

// TestPaddedPrimitiveTrailingCompletionPadsFillValueSlot verifies each atomic
// wrapper is followed by enough padding to complete its value slot.
func TestPaddedPrimitiveTrailingCompletionPadsFillValueSlot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		size      uintptr
		valueEnd  uintptr
		trailWant uintptr
	}{
		{
			name:      "PaddedUint64",
			size:      unsafe.Sizeof(PaddedUint64{}),
			valueEnd:  unsafe.Offsetof(PaddedUint64{}.value) + unsafe.Sizeof(PaddedUint64{}.value),
			trailWant: uintptr(CacheLinePadSize - atomicUint64Size),
		},
		{
			name:      "PaddedUint32",
			size:      unsafe.Sizeof(PaddedUint32{}),
			valueEnd:  unsafe.Offsetof(PaddedUint32{}.value) + unsafe.Sizeof(PaddedUint32{}.value),
			trailWant: uintptr(CacheLinePadSize - atomicUint32Size),
		},
		{
			name:      "PaddedInt64",
			size:      unsafe.Sizeof(PaddedInt64{}),
			valueEnd:  unsafe.Offsetof(PaddedInt64{}.value) + unsafe.Sizeof(PaddedInt64{}.value),
			trailWant: uintptr(CacheLinePadSize - atomicInt64Size),
		},
		{
			name:      "PaddedInt32",
			size:      unsafe.Sizeof(PaddedInt32{}),
			valueEnd:  unsafe.Offsetof(PaddedInt32{}.value) + unsafe.Sizeof(PaddedInt32{}.value),
			trailWant: uintptr(CacheLinePadSize - atomicInt32Size),
		},
		{
			name: "PaddedPointer",
			size: unsafe.Sizeof(PaddedPointer[padTestValue]{}),
			valueEnd: unsafe.Offsetof(PaddedPointer[padTestValue]{}.value) +
				unsafe.Sizeof(PaddedPointer[padTestValue]{}.value),
			trailWant: uintptr(CacheLinePadSize - atomicPointerSize),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			trailing := tc.size - tc.valueEnd
			if trailing < tc.trailWant {
				t.Fatalf("%s trailing completion pad = %d, want at least %d", tc.name, trailing, tc.trailWant)
			}
		})
	}
}

// TestPaddedPrimitiveAlignment verifies padded primitive sizes remain alignment-safe.
func TestPaddedPrimitiveAlignment(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		size      uintptr
		alignment uintptr
	}{
		{
			name:      "PaddedUint64",
			size:      unsafe.Sizeof(PaddedUint64{}),
			alignment: unsafe.Alignof(PaddedUint64{}),
		},
		{
			name:      "PaddedUint32",
			size:      unsafe.Sizeof(PaddedUint32{}),
			alignment: unsafe.Alignof(PaddedUint32{}),
		},
		{
			name:      "PaddedInt64",
			size:      unsafe.Sizeof(PaddedInt64{}),
			alignment: unsafe.Alignof(PaddedInt64{}),
		},
		{
			name:      "PaddedInt32",
			size:      unsafe.Sizeof(PaddedInt32{}),
			alignment: unsafe.Alignof(PaddedInt32{}),
		},
		{
			name:      "PaddedPointer",
			size:      unsafe.Sizeof(PaddedPointer[padTestValue]{}),
			alignment: unsafe.Alignof(PaddedPointer[padTestValue]{}),
		},
	}

	for _, tc := range tests {

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.alignment == 0 {
				t.Fatalf("unsafe.Alignof(%s{}) = 0, want non-zero alignment", tc.name)
			}

			if tc.size%tc.alignment != 0 {
				t.Fatalf(
					"unsafe.Sizeof(%s{}) = %d is not aligned to unsafe.Alignof(%s{}) = %d",
					tc.name,
					tc.size,
					tc.name,
					tc.alignment,
				)
			}
		})
	}
}

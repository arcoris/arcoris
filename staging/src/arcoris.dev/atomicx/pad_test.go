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
	"testing"
	"unsafe"
)

// TestCacheLinePadSize verifies the package's fixed false-sharing pad width.
func TestCacheLinePadSize(t *testing.T) {
	t.Parallel()

	if CacheLinePadSize != 128 {
		t.Fatalf("CacheLinePadSize = %d, want 128", CacheLinePadSize)
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

// TestPaddedPrimitiveSizesIncludeBothPads verifies primitive layouts have enough
// total space for both pads and the atomic wrapper they isolate.
func TestPaddedPrimitiveSizesIncludeBothPads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size uintptr
		min  uintptr
	}{
		{
			name: "PaddedUint64",
			size: unsafe.Sizeof(PaddedUint64{}),
			min:  uintptr(2*CacheLinePadSize) + unsafe.Sizeof(PaddedUint64{}.value),
		},
		{
			name: "PaddedUint32",
			size: unsafe.Sizeof(PaddedUint32{}),
			min:  uintptr(2*CacheLinePadSize) + unsafe.Sizeof(PaddedUint32{}.value),
		},
		{
			name: "PaddedInt64",
			size: unsafe.Sizeof(PaddedInt64{}),
			min:  uintptr(2*CacheLinePadSize) + unsafe.Sizeof(PaddedInt64{}.value),
		},
		{
			name: "PaddedInt32",
			size: unsafe.Sizeof(PaddedInt32{}),
			min:  uintptr(2*CacheLinePadSize) + unsafe.Sizeof(PaddedInt32{}.value),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.size < tt.min {
				t.Fatalf("unsafe.Sizeof(%s{}) = %d, want at least %d", tt.name, tt.size, tt.min)
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
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.offset < uintptr(CacheLinePadSize) {
				t.Fatalf("%s.value offset = %d, want at least %d", tt.name, tt.offset, CacheLinePadSize)
			}
		})
	}
}

// TestPaddedPrimitiveValueOffsetsKeepTrailingPad verifies each atomic wrapper is
// followed by the trailing pad. This protects arrays, slices, and following
// struct fields from sharing the atomic value's cache line.
func TestPaddedPrimitiveValueOffsetsKeepTrailingPad(t *testing.T) {
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
			trailWant: uintptr(CacheLinePadSize),
		},
		{
			name:      "PaddedUint32",
			size:      unsafe.Sizeof(PaddedUint32{}),
			valueEnd:  unsafe.Offsetof(PaddedUint32{}.value) + unsafe.Sizeof(PaddedUint32{}.value),
			trailWant: uintptr(CacheLinePadSize),
		},
		{
			name:      "PaddedInt64",
			size:      unsafe.Sizeof(PaddedInt64{}),
			valueEnd:  unsafe.Offsetof(PaddedInt64{}.value) + unsafe.Sizeof(PaddedInt64{}.value),
			trailWant: uintptr(CacheLinePadSize),
		},
		{
			name:      "PaddedInt32",
			size:      unsafe.Sizeof(PaddedInt32{}),
			valueEnd:  unsafe.Offsetof(PaddedInt32{}.value) + unsafe.Sizeof(PaddedInt32{}.value),
			trailWant: uintptr(CacheLinePadSize),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			trailingPad := tt.size - tt.valueEnd
			if trailingPad < tt.trailWant {
				t.Fatalf("%s trailing pad = %d, want at least %d", tt.name, trailingPad, tt.trailWant)
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
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.alignment == 0 {
				t.Fatalf("unsafe.Alignof(%s{}) = 0, want non-zero alignment", tt.name)
			}

			if tt.size%tt.alignment != 0 {
				t.Fatalf(
					"unsafe.Sizeof(%s{}) = %d is not aligned to unsafe.Alignof(%s{}) = %d",
					tt.name,
					tt.size,
					tt.name,
					tt.alignment,
				)
			}
		})
	}
}

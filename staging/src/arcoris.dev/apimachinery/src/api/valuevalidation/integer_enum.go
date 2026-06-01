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

package valuevalidation

import "arcoris.dev/apimachinery/api/fieldpath"

// validateSignedEnum checks signed integer enum membership when an enum exists.
func (v *validator) validateSignedEnum(path fieldpath.Path, got int64, enum []int64) {
	if len(enum) > 0 && !containsInt64(enum, got) {
		v.add(
			path,
			ErrEnumMismatch,
			ErrorReasonEnumMismatch,
			"integer value is not in enum",
		)
	}
}

// validateUnsignedEnum checks unsigned integer enum membership when an enum exists.
func (v *validator) validateUnsignedEnum(path fieldpath.Path, got uint64, enum []uint64) {
	if len(enum) > 0 && !containsUint64(enum, got) {
		v.add(
			path,
			ErrEnumMismatch,
			ErrorReasonEnumMismatch,
			"integer value is not in enum",
		)
	}
}

// signedEnum adapts smaller signed enum slices to int64.
func signedEnum[T ~int8 | ~int16 | ~int32](values []T) []int64 {
	out := make([]int64, len(values))
	for i, enumValue := range values {
		out[i] = int64(enumValue)
	}

	return out
}

// unsignedEnum adapts smaller unsigned enum slices to uint64.
func unsignedEnum[T ~uint8 | ~uint16 | ~uint32](values []T) []uint64 {
	out := make([]uint64, len(values))
	for i, enumValue := range values {
		out[i] = uint64(enumValue)
	}

	return out
}

// containsInt64 reports whether values contains target.
func containsInt64(values []int64, target int64) bool {
	for _, candidate := range values {
		if candidate == target {
			return true
		}
	}

	return false
}

// containsUint64 reports whether values contains target.
func containsUint64(values []uint64, target uint64) bool {
	for _, candidate := range values {
		if candidate == target {
			return true
		}
	}

	return false
}

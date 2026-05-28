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

package types

// ordered is the private comparison set used by descriptor-limit helpers.
//
// It is intentionally private. Public descriptor concepts remain exact:
// Int8View returns int8, Uint64View returns uint64, and every TypeCode still
// owns an exact payload slot.
type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// limit stores an optional descriptor rule without pointer allocation.
//
// The set bit distinguishes "unset" from an explicit zero value. Payloads use
// this private generic helper for low-level storage only; exact payload structs
// remain the descriptor boundary.
type limit[T any] struct {
	// value stores the configured rule value when set is true.
	value T
	// set records whether value was explicitly configured.
	set bool
}

// invalidRange reports whether min and max are both set and inverted.
func invalidRange[T ordered](min, max limit[T]) bool {
	return min.set && max.set && min.value > max.value
}

// enumBelowMin reports whether any enum value is below a configured minimum.
func enumBelowMin[T ordered](values []T, min limit[T]) bool {
	if !min.set {
		return false
	}
	for _, value := range values {
		if value < min.value {
			return true
		}
	}
	return false
}

// enumAboveMax reports whether any enum value is above a configured maximum.
func enumAboveMax[T ordered](values []T, max limit[T]) bool {
	if !max.set {
		return false
	}
	for _, value := range values {
		if value > max.value {
			return true
		}
	}
	return false
}

// validateLengthLimits checks non-negative size limits and ordering.
func validateLengthLimits(min, max limit[int], path string) error {
	if min.set && min.value < 0 {
		return typeError(path+".min", ErrInvalidType)
	}
	if max.set && max.value < 0 {
		return typeError(path+".max", ErrInvalidType)
	}
	if invalidRange(min, max) {
		return typeError(path, ErrInvalidType)
	}
	return nil
}

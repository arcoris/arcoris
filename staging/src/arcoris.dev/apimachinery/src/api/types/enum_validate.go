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

import "fmt"

// hasDuplicates reports whether values contains any repeated comparable value.
func hasDuplicates[T comparable](values []T) bool {
	_, _, ok := firstDuplicate(values)

	return ok
}

// firstDuplicate returns the first repeated enum value and its duplicate index.
func firstDuplicate[T comparable](values []T) (int, T, bool) {
	seen := make(map[T]struct{}, len(values))

	for i, value := range values {
		if _, ok := seen[value]; ok {
			return i, value, true
		}

		seen[value] = struct{}{}
	}

	var zero T

	return 0, zero, false
}

// enumBelowMin reports whether any enum value is below a configured minimum.
func enumBelowMin[T ordered](values []T, min limit[T]) bool {
	_, _, ok := firstEnumBelowMin(values, min)

	return ok
}

// firstEnumBelowMin returns the first enum value below a configured minimum.
func firstEnumBelowMin[T ordered](values []T, min limit[T]) (int, T, bool) {
	if !min.set {
		var zero T
		return 0, zero, false
	}

	for i, value := range values {
		if value < min.value {
			return i, value, true
		}
	}

	var zero T

	return 0, zero, false
}

// enumAboveMax reports whether any enum value is above a configured maximum.
func enumAboveMax[T ordered](values []T, max limit[T]) bool {
	_, _, ok := firstEnumAboveMax(values, max)

	return ok
}

// firstEnumAboveMax returns the first enum value above a configured maximum.
func firstEnumAboveMax[T ordered](values []T, max limit[T]) (int, T, bool) {
	if !max.set {
		var zero T
		return 0, zero, false
	}

	for i, value := range values {
		if value > max.value {
			return i, value, true
		}
	}

	var zero T

	return 0, zero, false
}

// validateEnumRules checks enum uniqueness and min/max membership.
func validateEnumRules[T ordered](path, descriptor string, values []T, min, max limit[T]) error {
	if index, value, ok := firstDuplicate(values); ok {
		return descriptorErrorf(
			path+".enum",
			ErrInvalidDescriptor,
			DescriptorErrorReasonDuplicateEnum,
			"%s enum values must be unique; duplicate value %v at index %d",
			descriptor,
			value,
			index,
		)
	}

	if index, value, ok := firstEnumBelowMin(values, min); ok {
		return descriptorErrorf(
			fmt.Sprintf("%s.enum[%d]", path, index),
			ErrInvalidDescriptor,
			DescriptorErrorReasonEnumBelowMinimum,
			"%s enum value %v is below minimum %v",
			descriptor,
			value,
			min.value,
		)
	}

	if index, value, ok := firstEnumAboveMax(values, max); ok {
		return descriptorErrorf(
			fmt.Sprintf("%s.enum[%d]", path, index),
			ErrInvalidDescriptor,
			DescriptorErrorReasonEnumAboveMaximum,
			"%s enum value %v is above maximum %v",
			descriptor,
			value,
			max.value,
		)
	}

	return nil
}
